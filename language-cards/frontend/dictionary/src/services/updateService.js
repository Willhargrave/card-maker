import axios from 'axios'

const updateWord = async (originalWord, newWord, newDefinition, userID, setId) => {
    try {
        const response = await axios.put('http://localhost:8080/update-word', {
            originalWord,
            word: newWord,
            definition: newDefinition,
            userID,
            setId
        });
        return response.data;
    } catch (error) {
        console.error('Error updating word', error);
        throw error;
    }
};
export default updateWord;