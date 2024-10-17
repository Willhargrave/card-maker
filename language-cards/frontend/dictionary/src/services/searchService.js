import axios from 'axios'

const search = async (word) => {
    try {
        const url = `http://localhost:8080/search?word=${word}`;
        const response = await axios.get(url);
        return response.data;
    } catch (error) {
        throw error;
    }
};

export default search;