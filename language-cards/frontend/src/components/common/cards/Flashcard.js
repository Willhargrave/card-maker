import React, {useState} from "react";
import UpdateWord from "../forms/wrappers/UpdateWord";
import AddWord from "../forms/wrappers/AddWord";
import Modal from "./Modal";

function Flashcard({word, definition, showDefinition, onShowDefinitionClick, onDelete, onCorrect, onIncorrect, showNext, setShowNext, onUpdate, onAdd, unseenWords, currentUser, currentSet, setCurrentSet}) {

    const [isAddModalOpen, setIsAddModalOpen] = useState(false)
    const [isUpdateModalOpen, setIsUpdateModalOpen] = useState(false)

   
    const handleDeleteClick = () => {
        onDelete(word)
    }

    return (
        <div>
            {unseenWords.length === 0 ? (
                <div className="text-center space-y-6">
                    <p className="text-xl">Well done! You've finished for today!</p>
                    <button className="bg-blue-500 text-white px-6 py-3 rounded-lg hover:bg-blue-600 transition-colors" 
                        onClick={() => setIsAddModalOpen(true)}>
                        Add New Word
                    </button>
                    <Modal isOpen={isAddModalOpen} onClose={() => setIsAddModalOpen(false)}>
                    <AddWord 
                    onAdd={onAdd} 
                    onClose={() => setIsAddModalOpen(false)}
                    currentUser={currentUser}
                    currentSet={currentSet}
                    /> 
                    </Modal>
                </div>
            ) : (
                <div className="space-y-12">
                    
                    <div className="flex justify-center gap-4">
                        <button className="bg-blue-500 text-white px-6 py-3 rounded-lg hover:bg-blue-600 transition-colors" 
                            onClick={() => setIsAddModalOpen(true)}>Add</button>
                        <button className="bg-green-500 text-white px-6 py-3 rounded-lg hover:bg-green-600 transition-colors" 
                            onClick={() => setIsUpdateModalOpen(true)}>Update</button>
                        <button className="bg-red-500 text-white px-6 py-3 rounded-lg hover:bg-red-600 transition-colors" 
                            onClick={handleDeleteClick}>Delete</button>
                        <button className="bg-gray-500 text-white px-6 py-3 rounded-lg hover:bg-gray-600 transition-colors" 
                            onClick={() => setCurrentSet(false)}>Return</button>
                    </div>
   
                    <Modal isOpen={isAddModalOpen} onClose={() => setIsAddModalOpen(false)}>
                    <AddWord 
                    onAdd={onAdd} 
                    onClose={() => setIsAddModalOpen(false)}
                    currentUser={currentUser}
                    currentSet={currentSet}
                    /> 
        </Modal>
        <Modal isOpen={isUpdateModalOpen} onClose={() => setIsUpdateModalOpen(false)}>
            <UpdateWord 
            word={word}
            definition={definition}
            onUpdate={onUpdate}
            onClose={() => setIsUpdateModalOpen(false)}
            currentUser={currentUser}
            currentSet={currentSet}
            />
        </Modal>
    

                    <div className="text-center space-y-8">
                        <h2 className="text-8xl font-bold">{word}</h2>
                        
                        {showDefinition && (
                            <p className="text-2xl text-gray-600">{definition}</p>
                        )}
                        
                        <button 
                            className="bg-purple-500 text-white px-8 py-3 rounded-lg hover:bg-purple-600 transition-colors"
                            onClick={onShowDefinitionClick}
                        >
                            {showDefinition ? 'Hide Definition' : 'Show Definition'}
                        </button>
    
                  
                        {showNext && (
                            <div className="flex justify-center gap-4 mt-8">
                                <button 
                                    className="bg-green-500 text-white px-8 py-3 rounded-lg hover:bg-green-600 transition-colors"
                                    onClick={() => onCorrect(word)}
                                >
                                    Correct
                                </button>
                                <button 
                                    className="bg-red-500 text-white px-8 py-3 rounded-lg hover:bg-red-600 transition-colors"
                                    onClick={onIncorrect}
                                >
                                    Wrong
                                </button>
                            </div>
                        )}
                    </div>
                </div>
            )}
        </div>
    );
}

export default Flashcard;


