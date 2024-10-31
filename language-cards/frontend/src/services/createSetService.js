import axios from 'axios'

const createSet = async (userID, setName) => {
    console.log("Sending to create-set:", { userID, setName });
    try {
        const response = await axios.post('http://localhost:8080/create-set', {
            userID,
            setName
        });
        return response;
    } catch (error) {
        console.error('Error creating set', error);
        throw error;
    }
}


export default createSet;

