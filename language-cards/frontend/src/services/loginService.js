import axios from 'axios'

const login = async (username, password) => {
    try {
        const url = `http://localhost:8080/login`;
        const response = await axios.post(url, {
            username,
            password
        });
        return response;
    } catch (error) {
        console.error('Error during login', error);
        throw error;
    }
}

export default login;