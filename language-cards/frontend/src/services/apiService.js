import axios from 'axios'

const addWord = async (word, definition, userID, setId) => {
    try {
        const response = await axios.post('http://localhost:8080/add-word', {
            word,
            definition,
            UserID: userID,
            setId 
        });
        return response.data;
    } catch (error) {
        console.error('Request that failed:', {
            word,
            definition,
            UserID: userID,
            setId: setId
        });
        throw error;
    }
}

export default addWord;