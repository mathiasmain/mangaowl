import { useNavigate, useLocation } from "react-router-dom";

const HomeButton = () => {
  const navigate = useNavigate();
  const location = useLocation();

  const handleClick = () => {
    if (location.pathname === "/") {
      window.location.reload(); // Recarrega a página se já estiver na home
    } else {
      navigate("/"); // Navega para a home se estiver em outra página
    }
  };

  return (
    <header className="fixed top-0 left-0 w-full shadow-md flex items-center p-3 z-50">
      <button
        onClick={handleClick}
        className="flex items-center space-x-2 px-4 py-2  hover:bg-gray-200 rounded-lg transition"
      >
        <img src="icon.png" alt="Logo" className="h-8" /> 
        {/* Substitua por seu ícone ou texto */}
      </button>
    </header>
  );
};

export default HomeButton;