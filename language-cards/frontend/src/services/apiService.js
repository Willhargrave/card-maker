import axios from 'axios'

const addWord = async (word, definition, userID, setId) => {
    try {
        const response = await axios.post('http://localhost:8080/add-word', {
            Word: word,
            definition,
            UserID: userID,
            setId 
        });
        const payload = { Word: word, UserID: userID, setId };
        console.log(payload)
        return response.data;
    } catch (error) {
        console.error('Request that failed:', {
            Word: word,
            definition,
            UserID: userID,
            setId: setId
        });
        throw error;
    }
}

export default addWord;