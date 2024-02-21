import {doThing} from './foo.js';

const h1 = document.querySelector('h1');
const helloBtn = document.querySelector('#helloBtn');
helloBtn.addEventListener('click', () => {
    h1.textContent = 'Goodbye'
});
doThing();