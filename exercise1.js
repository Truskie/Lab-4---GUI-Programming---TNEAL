// exercise1.js
const accordion = document.querySelector('.accordion');
accordion.addEventListener('click', function(event) {
 // 1. Find the closest trigger from event.target
 const trigger = event.target.closest('.accordion-item');
 // If null, the click was not on a trigger - return early.
 if(trigger == null) return;
 // 2. From the trigger, find the closest .accordion-item
 const currentItem = trigger.closest('.accordion-item');
 // 3. Close ALL accordion items (querySelectorAll + forEach + remove class)
 const allItems = document.querySelectorAll('.accordion-item');
 allItems.forEach(item => {
    if (item !== currentItem){
        item.classList.remove('open');
    }
 })
 // 4. Toggle .open on the clicked item only
 currentItem.classList.toggle('open');
});