// Простая hash-навигация между списком и деталями спектакля
function navigate() {
	const hash = window.location.hash || '#list';
	if (hash.startsWith('#play/')) {
		const playId = hash.replace('#play/', '');
		showDetailView(playId);
	} else {
		showListView();
	}
}

function showSection(sectionIdToShow) {
	const sections = ['view-list', 'view-detail'];
	sections.forEach((id) => {
		const el = document.getElementById(id);
		if (!el) return;
		el.style.display = id === sectionIdToShow ? 'block' : 'none';
	});
}

// Загрузка спектаклей с API и отрисовка списка
async function loadPlays() {
	try {
		const response = await fetch('/api/plays');
		if (!response.ok) throw new Error('HTTP ' + response.status);
		const plays = await response.json(); // Swagger: ответ — массив
		if (!Array.isArray(plays)) {
			throw new Error('Unexpected response format: expected array');
		}

		const playsContainer = document.getElementById('plays-list');
		if (!playsContainer) return;

		playsContainer.innerHTML = '';

		plays.forEach((play) => {
			const li = document.createElement('li');
			li.className = 'card';
			const badges = renderPerformanceBadges(play.performances);
			const thumb = play.poster_url ? `<img class="thumb" alt="${escapeAttr(play.title || '')}" src="${escapeAttr(play.poster_url)}"/>` : '<div class="thumb placeholder"></div>';
			const description = escapeHtml(play.description || '');
			li.innerHTML = `
				<div class="card-left">${thumb}</div>
				<div class="card-center">
					<a href="#play/${play.id}" class="card-title">${escapeHtml(play.title || '')}</a>
					<p class="card-desc">${description}</p>
					<div class="badges">${badges}</div>
				</div>
				<div class="card-right">
					<span class="author">${escapeHtml(play.author || '')}</span>
				</div>
			`;
			playsContainer.appendChild(li);
		});
	} catch (error) {
		console.error('Ошибка загрузки спектаклей:', error);
		const playsContainer = document.getElementById('plays-list');
		if (playsContainer) {
			playsContainer.innerHTML = '<li style="color: red;">Ошибка загрузки данных</li>';
		}
	}
}

function showListView() {
	showSection('view-list');
	// Если список ещё не подгружен (или хотим обновить), загружаем
	loadPlays();
}

async function showDetailView(playId) {
	showSection('view-detail');
	const container = document.getElementById('play-detail');
	if (!container) return;
	container.innerHTML = 'Загрузка...';
	try {
		const response = await fetch(`/api/plays/${encodeURIComponent(playId)}`);
		if (!response.ok) throw new Error('Не удалось получить данные');
		const play = await response.json();
		container.innerHTML = renderPlayDetailHtml(play);
	} catch (e) {
		container.innerHTML = '<div style="color: red;">Не удалось загрузить спектакль</div>';
	}
}

function renderPlayDetailHtml(play) {
	const meta = [
		`Жанр: ${escapeHtml(play.genre || '')}`,
		`Автор: ${escapeHtml(play.author || '')}`,
		`Длительность: ${typeof play.duration === 'number' ? play.duration + ' мин' : ''}`,
	];
	const poster = play.poster_url ? `<img alt="poster" src="${escapeAttr(play.poster_url)}" />` : '';
	return `
		<button type="button" class="btn btn-back" onclick="window.location.hash = '#list'">← Назад</button>
		<h2>${escapeHtml(play.title || '')}</h2>
		<div class="play-detail">
			<div class="play-poster">${poster}</div>
			<div class="play-meta">
				<p>${meta.filter(Boolean).join('<br>')}</p>
				<p>${escapeHtml(play.description || '')}</p>
			</div>
		</div>
		${renderPerformances(play.performances)}
	`;
}

function renderPerformances(performances) {
	if (!Array.isArray(performances) || performances.length === 0) return '';
	const rows = performances.map((p) => `
		<li>
			<span>${escapeHtml(p.date || '')}</span>
			<span class="status">${escapeHtml(p.status || '')}</span>
		</li>
	`).join('');
	return `
		<h3>Ближайшие показы</h3>
		<ul class="performances">${rows}</ul>
	`;
}

function renderPerformanceBadges(performances) {
	if (!Array.isArray(performances) || performances.length === 0) return '';
	return performances.slice(0, 4).map((p) => {
		const text = formatDateTimeRu(p.date);
		return `<span class="badge">${text}</span>`;
	}).join('');
}

function formatDateTimeRu(input) {
	if (!input) return '';
	// Try parse: ISO string or other; fallback to raw
	const d = new Date(input);
	if (isNaN(d.getTime())) return escapeHtml(String(input));
	const dtf = new Intl.DateTimeFormat('ru-RU', {
		day: '2-digit', month: 'short', hour: '2-digit', minute: '2-digit'
	});
	return escapeHtml(dtf.format(d).replace('.', ''));
}

function escapeHtml(value) {
	return String(value)
		.replace(/&/g, '&amp;')
		.replace(/</g, '&lt;')
		.replace(/>/g, '&gt;')
		.replace(/"/g, '&quot;')
		.replace(/'/g, '&#39;');
}

function escapeAttr(value) {
	return escapeHtml(value).replace(/"/g, '&quot;');
}

// Инициализация
document.addEventListener('DOMContentLoaded', () => {
	navigate();
});

window.addEventListener('hashchange', navigate);