import React, { useState } from 'react';
import login from '../services/loginService'; // Update the path as per your project structure
import "../styles/form.css"
function Login({setIsLoggedIn, loginSuccess}) {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');

    const handleLogin = async (e) => {
        e.preventDefault();
        try {
            const response = await login(username, password);
            localStorage.setItem('token', response.data.token);
            setIsLoggedIn(true); // Update your logged-in state
            loginSuccess();
            // Inform the parent component or redirect
        } catch (error) {
            console.error('Login error:', error);
            // Handle login errors
        }
    };

    return (
        <div>
            <h2>Login</h2>
            <form onSubmit={handleLogin} className="login-form">

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
                <button type="submit">Login</button>
            </form>
            {error && <p style={{ color: 'red' }}>{error}</p>}
        </div>
    );
}

export default Login;
