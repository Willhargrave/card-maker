import React, { useState } from 'react';
import register from '../services/registerService'; // Update the path as per your project structure
import "../styles/form.css"
import login from '../services/loginService';
function Register({setIsLoggedIn, onLoginSuccess}) {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');
    const [success, setSuccess] = useState('');

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            // Perform the registration
            await register(username, password);

            // Automatically log in the user after successful registration
            const loginResponse = await login(username, password);
            localStorage.setItem('token', loginResponse.data.token);
            setIsLoggedIn(true); // Set the logged-in state to true
            onLoginSuccess(); // Invoke any additional login success actions
        } catch (error) {
            setError('Registration failed: ' + error.message); // Handle errors
        }
    };

    return (
        <div className='login-container'>
            <h2>Register</h2>
            <form onSubmit={handleSubmit} className="login-form">
                <div>
                    <label>Username:</label>
                    <input 
                        type="text" 
                        value={username} 
                        onChange={(e) => setUsername(e.target.value)} 
                        required 
                    />
                </div>
                <div>
                    <label>Password:</label>
                    <input 
                        type="password" 
                        value={password} 
                        onChange={(e) => setPassword(e.target.value)} 
                        required 
                    />
                </div>
                <button type="submit">Register</button>
            </form>
            {error && <p style={{ color: 'red' }}>{error}</p>}
            {success && <p style={{ color: 'green' }}>{success}</p>}
        </div>
    );
}

export default Register;
