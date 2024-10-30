import React, { useState } from 'react';
import login from '../services/loginService'; // Update the path as per your project structure

function Login({setIsLoggedIn, loginSuccess}) {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');

    const handleLogin = async (e) => {
        e.preventDefault();
        try {
            const response = await login(username, password);
            localStorage.setItem('token', response.data.token);
            setIsLoggedIn(true); 
            loginSuccess();
   
        } catch (error) {
            console.error('Login error:', error);
            
        }
    };

    return (
        <div className="max-w-[400px] mx-auto p-5 border border-gray-200 rounded-lg shadow-sm">
    <h2 className="text-center mb-5 text-gray-700">Login</h2>
    <form onSubmit={handleLogin} className="flex flex-col gap-4">
        <div>
            <label className="font-bold mb-1.5">Username:</label>
            <input 
                type="text"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                required
                className="w-full p-2.5 border border-gray-300 rounded text-base"
            />
        </div>
        <div>
            <label className="font-bold mb-1.5">Password:</label>
            <input 
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
                className="w-full p-2.5 border border-gray-300 rounded text-base"
            />
        </div>
        <button type="submit" className="bg-blue-500 hover:bg-blue-700 text-white p-2.5 rounded text-base cursor-pointer transition-colors duration-300">Login/Register</button>
    </form>
    {error && <p className="text-red-500 text-center">{error}</p>}
</div>
    )    
}


export default Login;
