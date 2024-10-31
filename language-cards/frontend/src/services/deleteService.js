import axios from 'axios'

const deleteWord = async (word, userID, setId) => {
    try {
        const url = `http://localhost:8080/delete`;
        const response = await axios.delete(url, {
            data: { word, UserID: userID, setId } // Sending data in the request body
        });
        return response.data;
    } catch (error) {
        console.error('Request that failed:', {
            word,
            UserID: userID,
            setId: setId});
        throw error;
    }
};

export default deleteWord;