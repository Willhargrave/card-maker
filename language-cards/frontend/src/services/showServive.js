import axios from 'axios'

const showAllWords = async (userID, setId) => {
    try {
        // Use template literals to insert variables into the URL
        const url = `http://localhost:8080/show-words?UserID=${userID}&setId=${setId}`;
        const response = await axios.get(url);
        return response.data;
    } catch (error) {
        throw error;
    }
};

export default showAllWords;