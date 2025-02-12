const form = document.getElementById('signUpForm')
const username = document.getElementById('username-input')
const email = document.getElementById('email-input')
const password = document.getElementById('password-input')
const password2 = document.getElementById('confirm-password-input')
const errorMessage = document.getElementById('error-message')

// export let correctUsername = username;
// export let correctPassword = password;

form.addEventListener('submit', (e) => {
    e.preventDefault();
    let errors = [];

    console.log("Checking email value: ", email.value)
    if(!email.value.trim() || !email.value.includes('@')){
    errors.push('Please enter a valid email address');
    email.parentElement.classList.add('incorrect');
    }

    if(!username.value.trim()){
    errors.push('Please enter a username');
    username.parentElement.classList.add('incorrect');
    }

    if(!password.value.trim()){
    errors.push('Please enter a password.');
    password.parentElement.classList.add('incorrect');
    }

    if(password.value !== password2.value){
    errors.push('Passwords do not match.');
    password2.parentElement.classList.add('incorrect');
    }
    else if(email.value.trim() && email.value.includes('@') && username.value.trim() && password.value.trim() && password.value === password2.value){
    form.submit();
    window.location.href = '/templates/login.html';
    }

    if(errors.length > 0){
    // If there are any errors
    e.preventDefault();
    errorMessage.innerText  = errors.join(". ");
    } else{
        form.submit();
        window.location.href = '/templates/login.html';
    }
});

const allInputs = [username, email, password, password2].filter(input => input != null);

allInputs.forEach(input => {
    input.addEventListener('input', () => {
        if(input.parentElement.classList.contains('incorrect')){
            input.parentElement.classList.remove('incorrect');
            errorMessage.innerText = '';
        }
    })
})

//Fixes that need to be added:
// 1. Find a way to check if the email is already in use
// 2. Add a way to check if the username is already in use
// 3. Add a way to check if the password is strong enough
// 4. Compare username and passwords to a database to see if they're actually correct, will probably have to combine signUp.js & login.js
