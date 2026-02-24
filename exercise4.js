// exercise4.js
const editList = document.querySelector('#edit-list');
editList.addEventListener('dblclick', function(event) {
 // 1. Find closest .edit-item from event.target; return if null
 const item = event.target.closest('.edit-item');
 if(!item) return;
 // 2. Return early if item already has .editing class
 if(item.classList.contains('editing')) return;
 // 3. Save original text, clear item, create and append input
 const originalText = item.textContent;
 item.textContent = '';
 const input = document.createElement('input');
 input.value = originalText;
 item.appendChild(input);
 item.classList.add('editing');
 input.focus();
 // -- Helper: commit the edit --
 let processing = false;

 function commitEdit() {
    if(processing) return;
    processing = true;

    const newText = input.value.trim() || originalText;

 // 4. Set item.textContent to newText, remove .editing
 item.textContent = newText;
 item.classList.remove('editing');
 }
 // -- Helper: cancel the edit --
 function cancelEdit() {
    if (processing) return;
    processing = true;

 // 5. Restore originalText, remove .editing
 item.textContent = originalText;
 item.classList.remove('editing');
 }
 // 6. Listen for 'keydown' on input
 // Enter -> commitEdit()
 // Escape -> cancelEdit()
input.addEventListener('keydown', function(e)
{
    if (e.key === 'Enter')
    {
        e.preventDefault();
        commitEdit();
    }
    else if (e.key === 'Escape') 
    {
        e.preventDefault();
        cancelEdit();
    }
});
 // 7. Listen for 'blur' on input -> commitEdit()
input.addEventListener('blur', function()
{
    if(!processing)
    {
        commitEdit();
    }
})
});
