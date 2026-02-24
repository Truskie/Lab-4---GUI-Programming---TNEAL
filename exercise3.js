// exercise3.js
const filterBar = document.querySelector('#filter-bar');
const cards = document.querySelectorAll('.card');
const filterBtns = document.querySelectorAll('.filter-btn');
filterBar.addEventListener('click', function(event) {
 // 1. Guard: if event.target does not match '.filter-btn', return
 if(!event.target.matches('.filter-btn')) return;
 // 2. Update active class on all buttons
 const filterButtons = document.querySelectorAll('.filter-btn');
 filterButtons.forEach(btn => {
    btn.classList.remove('active');
 })
 event.target.classList.add('active');
 // 3. Read the filter value from event.target.dataset.filter
 const filterVal = event.target.dataset.filter;
 // 4. Loop through cards
 cards.forEach(card => {
     // - If filter === 'all': remove .hidden from every card
    if(filterVal === 'all')
        {
            card.classList.remove('hidden');
        }
    else
        // - Otherwise: compare card.dataset.category to filter
        // Add .hidden if no match, remove .hidden if match
    {
        const hide = card.dataset.category !==filterVal;
        card.classList.toggle('hidden', hide);
    }
 });
});