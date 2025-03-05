import {
  BrowserRouter as Router,
  Routes,
  Route,
} from "react-router-dom";
import LoginPage from './pages/Login';
import HomePage from "./pages/Home";
import "@tabler/core/dist/css/tabler.min.css";
import 'bootstrap/dist/js/bootstrap.js';
import AuthRequired from "./pages/AuthRequired";


export default function App() {
  return (
    <Router>
      <Routes>
        <Route path='/' element={<AuthRequired><HomePage /></AuthRequired>} />
      </Routes>
      <Routes>
        <Route path='/login/' element={<LoginPage />} />
      </Routes>
    </Router>
  );
}
