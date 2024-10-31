import React, { useState } from 'react';
import login from '../../../services/loginService';
import register from '../../../services/registerService';

function AuthForm({ setIsLoggedIn, onSuccess }) {
    const [isRegistering, setIsRegistering] = useState(false);
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            if (isRegistering) {
                await register(username, password);
                const loginResponse = await login(username, password);
                localStorage.setItem('token', loginResponse.data.token);
            } else {
                const response = await login(username, password);
                localStorage.setItem('token', response.data.token);
            }
            setIsLoggedIn(true);
            onSuccess();
        } catch (error) {
            setError(`${isRegistering ? 'Registration' : 'Login'} failed: ${error.message}`);
        }
    };

    return (
        <div className="max-w-[400px] mx-auto p-5 border border-gray-200 rounded-lg shadow-sm">
            <h2 className="text-center mb-5 text-gray-700">
                {isRegistering ? 'Register' : 'Login'}
            </h2>
            <form onSubmit={handleSubmit} className="flex flex-col gap-4">
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
                <button 
                    type="submit" 
                    className="bg-blue-500 hover:bg-blue-700 text-white p-2.5 rounded text-base cursor-pointer transition-colors duration-300"
                >
                    {isRegistering ? 'Register' : 'Login'}
                </button>
            </form>
            <button 
                onClick={() => setIsRegistering(!isRegistering)}
                className="w-full mt-4 text-blue-500 hover:text-blue-700"
            >
                {isRegistering ? 'Already have an account? Login' : 'Need an account? Register'}
            </button>
            {error && <p className="text-red-500 text-center mt-4">{error}</p>}
        </div>
    );
}

export default AuthForm;