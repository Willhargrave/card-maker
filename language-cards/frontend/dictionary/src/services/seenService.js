import axios from 'axios'

const updateSeen = async (word, userID, setId) => {
    try {
        const response = await axios.post('http://localhost:8080/seen', {
            word, userID, setId
        });
        return response.data;
    } catch (error) {
        console.error('Error updating seen status', error);
        throw error;
    }
};


export default updateSeen;