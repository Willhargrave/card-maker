import axios from 'axios'

const addWord = async (word, definition, userID, setId) => {
    try {
        const response = await axios.post('http://localhost:8080/add-word', {
            word,
            definition,
            userID,
            setId 
        });
        return response.data;
    } catch (error) {
        console.error('Error adding word', error);
        throw error;
    }
}

export default addWord;