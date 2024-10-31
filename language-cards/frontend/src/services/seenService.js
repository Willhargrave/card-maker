import axios from 'axios'

const updateSeen = async (word, userID, setId) => {
    try {
        const response = await axios.post('http://localhost:8080/seen', {
            word, 
            UserID: userID, 
            setId
        });
        return response.data;
    } catch (error) {
        console.error('Request that failed:', {
            word,
            UserID: userID,
            setId: setId
        });
        throw error;
    }
};


export default updateSeen;