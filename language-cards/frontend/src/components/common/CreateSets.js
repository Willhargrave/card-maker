import React, {useState, useEffect} from "react";
import createSet from "../../services/createSetService";

function CreateSets({currentUser, onClose}) {
   const [setName, setSetName] = useState('');
   const [message, setMessage] = useState('');
   


   const handleSubmit = async (e) => {
    e.preventDefault();
    if (currentUser && currentUser.userID) {
    try {
        const response = await createSet(currentUser.userID, setName);
        setMessage(`Set "${setName}" created successfully!`);
    } catch (error) {
        setMessage('Error creating set');
    }} 
};

   return (
    <div>
        <h2>Create New Set</h2>
        <form onSubmit={handleSubmit}>
            <div>
                <label>Set Name:</label>
                <input 
                    type="text" 
                    value={setName} 
                    onChange={(e) => setSetName(e.target.value)} 
                    required 
                />
            </div>
            <button type="submit">Create Set</button>
        </form>
        {message && <p>{message}</p>}
    </div>
);
}

export default CreateSets;