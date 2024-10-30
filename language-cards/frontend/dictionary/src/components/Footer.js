import React from 'react';


const Footer = () => {
    return (
        <footer className="text-center p-5 bg-white/50 text-gray-700 fixed left-0 bottom-0 w-full">
            <p>© {new Date().getFullYear()} Learnify</p>
        </footer>
    );
};

export default Footer;