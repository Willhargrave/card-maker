import React from 'react';
import '../styles/footer.css'; // Import the CSS file

const Footer = () => {
    return (
        <footer className="footer">
            <p>© {new Date().getFullYear()} Learnify</p>
        </footer>
    );
};

export default Footer;