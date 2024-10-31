import axios from 'axios';

const fetchUserSets = async (userID) => {
    const token = localStorage.getItem('token');
    try {
        const url = `http://localhost:8080/fetch-sets?UserID=${userID}`;
        const response = await axios.get(url, {
            headers: {
                Authorization: `Bearer ${token}`
            }
        });
        return response.data;
    } catch (error) {
        console.error('Error fetching user sets', error);
        throw error;
    }
};


export default fetchUserSets;