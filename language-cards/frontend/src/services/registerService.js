import axios from 'axios'


const register = async (username, password) => {
    try {
        const url = `http://localhost:8080/register`;
        const response = await axios.post(url, {
            username,
            password
        });
        return response.data;
    } catch (error) {
        console.error('Error during registration', error);
        throw error;
    }
}

export default register;