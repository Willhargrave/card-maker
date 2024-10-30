import React from 'react';


function Header({ isLoggedIn, handleLogout, setIsCreateModalOpen }) {
    return (
        <div className="bg-green-500 text-white flex justify-between items-center p-2.5 mb-7">
        <div className="flex-grow text-center">
            <h1 className="text-2xl font-bold [text-shadow:2px_2px_4px_rgba(0,0,0,0.5)]">Learnify</h1>
        </div>
        {isLoggedIn && (
            <button className="bg-transparent border-none text-inherit cursor-pointer text-right p-2.5 text-base hover:underline hover:decoration-white" onClick={handleLogout}>
                Logout
            </button>
        )}
        <button className="bg-transparent border-none text-inherit cursor-pointer text-right p-2.5 text-base hover:underline hover:decoration-white" onClick={() => setIsCreateModalOpen(true)}>
            Add Set
        </button>
     </div>
     
    );
}



export default Header;