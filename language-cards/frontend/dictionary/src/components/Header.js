import React from 'react';
import '../styles/header.css'

function Header({ isLoggedIn, handleLogout, setIsCreateModalOpen }) {
    return (
        <div className="header">
            <div className="header-title">
                <h1>Learnify</h1>
            </div>
            {isLoggedIn && (
                <button className="logout-button" onClick={handleLogout}>
                    Logout
                </button>
            )}
            <button className="logout-button" onClick={() => setIsCreateModalOpen(true)}>Add Set</button>
        </div>
    );
}


const styles = {
}

export default Header;