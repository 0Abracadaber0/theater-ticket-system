// Простая hash-навигация между списком и деталями спектакля
function navigate() {
    const hash = window.location.hash || '#list';
    if (hash.startsWith('#play/')) {
        const playId = hash.replace('#play/', '');
        showDetailView(playId);
    } else if (hash.startsWith('#performance/')) {
        const performanceId = hash.replace('#performance/', '');
        showPerformanceSeatsView(performanceId);
    } else if (hash === '#bookings') {
        showBookingsView();
    } else {
        showListView();
    }
}

function showSection(sectionIdToShow) {
    const sections = ['view-list', 'view-detail', 'view-seats', 'view-bookings'];
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
        const plays = await response.json();
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

async function showPerformanceSeatsView(performanceId) {
    showSection('view-seats');
    const container = document.getElementById('seats-container');
    if (!container) return;
    container.innerHTML = 'Загрузка мест...';

    try {
        const [perfResp, seatsResp] = await Promise.all([
            fetch(`/api/performances/${encodeURIComponent(performanceId)}`),
            fetch(`/api/performances/${encodeURIComponent(performanceId)}/seats`)
        ]);

        if (!perfResp.ok || !seatsResp.ok) throw new Error('Ошибка загрузки данных');

        const performance = await perfResp.json();
        const seats = await seatsResp.json();

        container.innerHTML = renderSeatsSelectionHtml(performance, seats);
        initSeatsSelection(performanceId, seats);
    } catch (e) {
        console.error('Ошибка загрузки мест:', e);
        container.innerHTML = '<div style="color: red;">Не удалось загрузить места</div>';
    }
}

async function showBookingsView() {
    showSection('view-bookings');
    const container = document.getElementById('bookings-container');
    if (!container) return;

    const phone = prompt('Введите ваш номер телефона для просмотра бронирований:');
    if (!phone) {
        window.location.hash = '#list';
        return;
    }

    container.innerHTML = 'Загрузка бронирований...';

    try {
        const response = await fetch(`/api/bookings?phone=${encodeURIComponent(phone)}`);
        if (!response.ok) throw new Error('Ошибка загрузки бронирований');

        const bookings = await response.json();
        container.innerHTML = renderBookingsHtml(bookings);
    } catch (e) {
        console.error('Ошибка загрузки бронирований:', e);
        container.innerHTML = '<div style="color: red;">Не удалось загрузить бронирования</div>';
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
		${renderPerformances(play.performances, play.id)}
	`;
}

function renderPerformances(performances, playId) {
    if (!Array.isArray(performances) || performances.length === 0) return '';
    const rows = performances.map((p) => `
		<li>
			<div class="perf-info">
				<span class="perf-date">${formatDateTimeRu(p.date)}</span>
				<span class="status-badge status-${p.status}">${escapeHtml(p.status || '')}</span>
			</div>
			<button type="button" class="btn" onclick="window.location.hash='#performance/${p.id}'">
				Выбрать места
			</button>
		</li>
	`).join('');
    return `
		<h3 class="perf-heading">Ближайшие показы</h3>
		<ul class="performances">${rows}</ul>
	`;
}

function renderSeatsSelectionHtml(performance, seats) {
    const playTitle = performance.play ? escapeHtml(performance.play.title) : 'Спектакль';
    const perfDate = formatDateTimeRu(performance.date);

    const seatsByRow = {};
    seats.forEach(seat => {
        if (!seat.seat) return;
        const row = seat.seat.row;
        if (!seatsByRow[row]) seatsByRow[row] = [];
        seatsByRow[row].push(seat);
    });

    const rows = Object.keys(seatsByRow).sort((a, b) => a - b);
    const seatsHtml = rows.map(row => {
        const rowSeats = seatsByRow[row].sort((a, b) => a.seat.number - b.seat.number);
        const seatsInRow = rowSeats.map(seat => {
            const status = seat.status || 'available';
            const disabled = status !== 'available' ? 'disabled' : '';
            const category = seat.seat.category || '';
            return `
				<button type="button" 
					class="seat seat-${status} seat-${category}" 
					data-seat-id="${seat.id}"
					data-price="${seat.price}"
					${disabled}>
					${seat.seat.number}
				</button>
			`;
        }).join('');

        return `
			<div class="seat-row">
				<span class="row-label">Ряд ${row}</span>
				<div class="seats-in-row">${seatsInRow}</div>
			</div>
		`;
    }).join('');

    return `
		<button type="button" class="btn btn-back" onclick="history.back()">← Назад</button>
		<h2>${playTitle}</h2>
		<p class="perf-date-large">${perfDate}</p>
		
		<div class="legend">
			<div class="legend-item"><span class="legend-seat seat-available"></span> Свободно</div>
			<div class="legend-item"><span class="legend-seat seat-reserved"></span> Забронировано</div>
			<div class="legend-item"><span class="legend-seat seat-sold"></span> Продано</div>
		</div>
		
		<div class="seats-hall">
			${seatsHtml}
		</div>
		
		<div class="booking-panel" id="booking-panel" style="display: none;">
			<div class="booking-info">
				<p>Выбрано мест: <span id="selected-count">0</span></p>
				<p>Сумма: <span id="total-price">0</span> руб.</p>
			</div>
			<button type="button" class="btn btn-primary" onclick="proceedToBooking()">
				Забронировать
			</button>
		</div>
	`;
}

function renderBookingsHtml(bookings) {
    if (!bookings || bookings.length === 0) {
        return '<p>У вас пока нет бронирований</p>';
    }

    const rows = bookings.map(booking => {
        const playTitle = booking.performance && booking.performance.play
            ? escapeHtml(booking.performance.play.title)
            : 'Спектакль';
        const perfDate = booking.performance
            ? formatDateTimeRu(booking.performance.date)
            : '';
        const seatsCount = booking.seats ? booking.seats.length : 0;
        const statusClass = `status-${booking.status}`;

        return `
			<li class="booking-card">
				<div class="booking-header">
					<h3>${playTitle}</h3>
					<span class="status-badge ${statusClass}">${escapeHtml(booking.status)}</span>
				</div>
				<p class="booking-date">${perfDate}</p>
				<p>Мест: ${seatsCount}</p>
				<p class="booking-price">Сумма: ${booking.total_price} руб.</p>
				${booking.status === 'pending' ? `
					<button type="button" class="btn btn-cancel" onclick="cancelBooking('${booking.id}')">
						Отменить
					</button>
				` : ''}
			</li>
		`;
    }).join('');

    return `
		<button type="button" class="btn btn-back" onclick="window.location.hash = '#list'">← Назад</button>
		<h2>Мои бронирования</h2>
		<ul class="bookings-list">${rows}</ul>
	`;
}

let selectedSeats = [];
let currentPerformanceId = null;

function initSeatsSelection(performanceId, seats) {
    currentPerformanceId = performanceId;
    selectedSeats = [];

    document.querySelectorAll('.seat:not([disabled])').forEach(btn => {
        btn.addEventListener('click', function() {
            const seatId = this.getAttribute('data-seat-id');
            const price = parseInt(this.getAttribute('data-price') || '0');

            if (this.classList.contains('selected')) {
                this.classList.remove('selected');
                selectedSeats = selectedSeats.filter(s => s.id !== seatId);
            } else {
                this.classList.add('selected');
                selectedSeats.push({ id: seatId, price: price });
            }

            updateBookingPanel();
        });
    });
}

function updateBookingPanel() {
    const panel = document.getElementById('booking-panel');
    const countEl = document.getElementById('selected-count');
    const priceEl = document.getElementById('total-price');

    if (!panel || !countEl || !priceEl) return;

    const count = selectedSeats.length;
    const total = selectedSeats.reduce((sum, seat) => sum + seat.price, 0);

    countEl.textContent = count;
    priceEl.textContent = total;

    panel.style.display = count > 0 ? 'flex' : 'none';
}

async function proceedToBooking() {
    if (selectedSeats.length === 0) {
        alert('Выберите хотя бы одно место');
        return;
    }

    const phone = prompt('Введите ваш номер телефона:');
    if (!phone) return;

    const name = prompt('Введите ваше имя:');
    if (!name) return;

    try {
        const response = await fetch('/api/bookings', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                phone: phone,
                name: name,
                performance_id: currentPerformanceId,
                seat_ids: selectedSeats.map(s => s.id)
            })
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || 'Ошибка бронирования');
        }

        const booking = await response.json();
        alert('Бронирование успешно создано!\nНомер бронирования: ' + booking.id);
        window.location.hash = '#bookings';
    } catch (e) {
        console.error('Ошибка бронирования:', e);
        alert('Ошибка: ' + e.message);
    }
}

async function cancelBooking(bookingId) {
    if (!confirm('Вы уверены, что хотите отменить бронирование?')) return;

    try {
        const response = await fetch(`/api/bookings/${bookingId}/cancel`, {
            method: 'PATCH'
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || 'Ошибка отмены');
        }

        alert('Бронирование отменено');
        showBookingsView();
    } catch (e) {
        console.error('Ошибка отмены:', e);
        alert('Ошибка: ' + e.message);
    }
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