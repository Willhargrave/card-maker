function Modal({ isOpen, onClose, children }) {
    if (!isOpen) return null;

    return (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center">
            <div className="bg-white p-5 rounded-lg shadow-lg max-w-[500px] w-full">
                <button onClick={onClose} className="float-right border-none bg-none text-base cursor-pointer">X</button>
                {children}
            </div>
        </div>
    )
}

export default Modal;