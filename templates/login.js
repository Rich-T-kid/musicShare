const username = document.getElementById('username-input');
const password = document.getElementById('password-input');
const form = document.getElementById('loginForm');
const errorMessage = document.getElementById('error-message');
import { correctUsername, correctPassword } from './signUp.js';

form.addEventListener('submit', (e) => {
    e.preventDefault();
    let errors = [];

    if(!username.value.trim()){
        errors.push('Please enter a username');
        username.parentElement.classList.add('incorrect');
    } else if (username.value !== correctUsername.value){
        errors.push('Username does not exist');
        username.parentElement.classList.add('incorrect');
    }
    if(!password.value.trim()){
        errors.push('Please enter a password.');
        password.parentElement.classList.add('incorrect');
    } else if (password.value !== correctPassword.value){
        errors.push('Password is incorrect');
        password.parentElement.classList.add('incorrect');
    }

    if(username.value.trim() && password.value.trim()){
        form.submit();
        window.location.href = '/templates/SongofDay.html';
    }

    if(errors.length > 0){
        e.preventDefault();
        errorMessage.innerText = errors.join(". ");
    }
});

const allInputs = [username, password].filter(input => input != null);

allInputs.forEach(input => {
    input.addEventListener('input', () => {
        if(input.parentElement.classList.contains('incorrect')){
            input.parentElement.classList.remove('incorrect');
            errorMessage.innerText = '';
        }
    })
})
