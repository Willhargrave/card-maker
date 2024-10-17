import React, {useState} from 'react'
import search from '../services/searchService'

function Search() {
    const [word, setWord] = useState('')
    const [definition, setDefinition] = useState('')
    const [message, setMessage] = useState('')

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            const result = await search(word);
            setDefinition(result);
            setMessage('');
        } catch (error) {
            setMessage(`Sorry, ${word} is not in your deck`);
        }
    };

    return (
        <div>
            <form onSubmit={handleSubmit}>
                <input type="text" value={word} onChange={(e) => setWord(e.target.value)} placeholder="Word" required />
                <button type="submit">Search</button>
            </form>
            {definition && <p>Definition: {definition}</p>}
            {message && <p>{message}</p>}
        </div>
    );
}

export default Search;
  
  