import React, {useState} from "react";
import UpdateWord from "./UpdateWord";
import AddWord from "./AddWord";
import Modal from "./Modal";

function Flashcard({word, definition, showDefinition, onShowDefinitionClick, onDelete, onCorrect, onIncorrect, showNext, setShowNext, onUpdate, onAdd, unseenWords, currentUser, currentSet, setCurrentSet}) {

    const [isAddModalOpen, setIsAddModalOpen] = useState(false)
    const [isUpdateModalOpen, setIsUpdateModalOpen] = useState(false)

    console.log("Flashcard rendering with props:", { 
        word, 
        showDefinition, 
        unseenWords: unseenWords.length 
    });

    const handleDeleteClick = () => {
        onDelete(word)
    }

    return (
        <div>
        {unseenWords.length === 0 ? (
            <>
                <p>Well done! You've finished for today!</p>
                <button className="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600" onClick={() => setIsAddModalOpen(true)}>Add New Word</button>
                <Modal isOpen={isAddModalOpen} onClose={() => setIsAddModalOpen(false)}>
                    <AddWord onAdd={onAdd} currentUser={currentUser} currentSet={currentSet} onClose={() => setIsAddModalOpen(false)}/>
                </Modal>
            </>
        ) : (
            <>
                <div className="flex justify-around items-center my-2.5">
                    <button className="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600" onClick={() => setIsAddModalOpen(true)}>Add</button>
                    <button className="bg-green-500 text-white px-4 py-2 rounded hover:bg-green-600" onClick={() => setIsUpdateModalOpen(true)}>Update</button>
                    <button className="bg-red-500 text-white px-4 py-2 rounded hover:bg-red-600" onClick={handleDeleteClick}>Delete</button>
                    <button className="bg-gray-500 text-white px-4 py-2 rounded hover:bg-gray-600" onClick={() => setCurrentSet(false)}>Return</button>
                </div>
                <Modal isOpen={isAddModalOpen} onClose={() => setIsAddModalOpen(false)}>
                    <AddWord onAdd={onAdd} currentUser={currentUser} currentSet={currentSet} onClose={() => setIsAddModalOpen(false)}/>
                </Modal>
                <Modal isOpen={isUpdateModalOpen} onClose={() => setIsUpdateModalOpen(false)}>
                    <UpdateWord word={word} definition={definition} onUpdate={onUpdate} currentSet={currentSet} currentUser={currentUser} onClose={() => setIsUpdateModalOpen(false)} />
                </Modal>
     
                <div className="text-gray-700 mt-0 text-center text-2xl">
                    <h2 className="text-[84px]">{word}</h2>
                    {showDefinition && <p>{definition}</p>}
                    <button className="bg-purple-500 text-white px-4 py-2 rounded hover:bg-purple-600"
                        onClick={() => {
                            console.log("Show definition button clicked");
                            onShowDefinitionClick();
                        }}>
                        {showDefinition ? 'Hide Definition' : 'Show Definition'}
                    </button>
                    {showNext && (
                        <>
                            <button className="bg-green-500 text-white px-4 py-2 rounded hover:bg-green-600 ml-2" onClick={() => onCorrect(word)}>Correct</button>
                            <button className="bg-red-500 text-white px-4 py-2 rounded hover:bg-red-600 ml-2" onClick={onIncorrect}>Wrong</button>
                        </>
                    )}
                </div>
            </>
        )}
     </div>
    );
}

export default Flashcard;