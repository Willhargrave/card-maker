import axios from 'axios'

const deleteWord = async (word, userID, setId) => {
    try {
        const url = `http://localhost:8080/delete`;
        const payload = { word, UserID: userID, setId };
        console.log('Delete payload before sending:', JSON.stringify(payload, null, 2));
        
        const response = await axios.delete(url, {
            data: payload
        });
        return response.data;
    } catch (error) {
        console.error('Delete request failed:', {
            sentPayload: { word, UserID: userID, setId },
            error: error.response?.data
        });
        throw error;
    }
};

export default deleteWord;