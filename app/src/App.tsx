import {
  BrowserRouter as Router,
  Routes,
  Route,
} from "react-router-dom";
import LoginPage from './pages/Login';
import HomePage from "./pages/Home";
import "@tabler/core/dist/css/tabler.min.css";


export default function App() {
  return (
      <Router>
        <Routes>
          <Route path='/login/' element={<LoginPage />} />
          <Route path='/' element={<HomePage />} />
        </Routes>
      </Router>
  );
}
