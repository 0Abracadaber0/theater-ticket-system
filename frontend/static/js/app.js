// Загрузка спектаклей с API
async function loadPlays() {
    try {
        const response = await fetch('/api/plays');
        const data = await response.json();

        const playsContainer = document.getElementById('plays-list');
        if (!playsContainer) return;

        playsContainer.innerHTML = '';

        data.plays.forEach(play => {
            const li = document.createElement('li');
            li.innerHTML = `
                <a href="/plays/${play.id}">${play.title}</a>
                <span>${play.author}</span>
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

// Загружаем данные при загрузке страницы
document.addEventListener('DOMContentLoaded', () => {
    if (document.getElementById('plays-list')) {
        loadPlays();
    }
});