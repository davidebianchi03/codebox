
import React from "react";

export default function HomePage() {

    // const [email, setEmail] = useState("");
    // const [password, setPassword] = useState("");
    // const [error, setError] = useState("");

    // const navigate = useNavigate();

    // const IsAuthenticated = useCallback(async () => {
    //     // redirect to home if user is already authenticated
    //     let [status, statusCode] = await Http.Request(`${Http.GetServerURL()}/api/v1/auth/user-details`, "GET", null);
    //     if (status === RequestStatus.OK && statusCode === 200) {
    //         navigate("/")
    //         return
    //     }
    // }, [navigate])

    // useEffect(() => {
    //     IsAuthenticated();
    // }, [IsAuthenticated]);

    // const SubmitLoginForm = async (event: any) => {
    //     event.preventDefault();

    //     // validate fields
    //     if (email === "" || password === "") {
    //         setError("Missing email or password");
    //         return;
    //     }

    //     // process login
    //     let [status, jwtToken, expirationDate] = await Http.Login(email, password);
    //     if (status === LoginStatus.OK) {
    //         setError("");
    //         document.cookie = `jwtToken=${jwtToken};expires=${expirationDate.toUTCString()};domain=${window.location.hostname}`;
    //         document.cookie = `jwtToken=${jwtToken};expires=${expirationDate.toUTCString()};domain=.${window.location.hostname}`;
    //         navigate("/")
    //         return
    //     } else {
    //         document.cookie = `jwtToken=invalidtoken;expires=Thu, 01 Jan 1970 00:00:01 GMT;domain=${window.location.hostname}`;
    //         document.cookie = `jwtToken=invalidtoken;expires=Thu, 01 Jan 1970 00:00:01 GMT;domain=.${window.location.hostname}`;
    //         if (status === LoginStatus.INVALID_CREDENTIALS) {
    //             setError("Invalid credentials");
    //         } else {
    //             setError("Unknown error, check that server is reachable");
    //         }
    //     }
    // }

    // return (
    //     <div style={{
    //         display: "flex",
    //         alignItems: "center",
    //         justifyContent: "center",
    //         width: "100%",
    //         height: "100%",
    //         background: "var(--background-color)"
    //     }}>
    //         <div style={{ width: "350px", display: "flex", flexDirection: "column" }}>
    //             <div style={{ display: "flex", justifyContent: "center", marginTop: "10pt", marginBottom: "20pt" }}>
    //                 <img src={CodeboxLogoWhite} style={{ maxWidth: "250px" }} alt="Codebox logo" />
    //             </div>
    //             <div style={{ textAlign: "center", marginBottom: "10pt", color: "var(--red)" }}>
    //                 {error}
    //             </div>
    //             <form onSubmit={SubmitLoginForm}>
    //                 <TextInput
    //                     label={"Email"}
    //                     placeholder={"john@doe.com"}
    //                     style={{ width: "calc(100% - 15pt)" }}
    //                     onTextChanged={(event) => setEmail(event.target.value)}
    //                     autocomplete="email"
    //                     name="email"
    //                 />
    //                 <TextInput
    //                     label={"Password"}
    //                     placeholder={"password"}
    //                     secure={true}
    //                     style={{ width: "calc(100% - 15pt)", marginTop: "10pt" }}
    //                     onTextChanged={(event) => setPassword(event.target.value)}
    //                     autocomplete="password"
    //                     name="password"
    //                 />
    //                 <Button
    //                     style={{
    //                         display: "flex",
    //                         justifyContent: "center",
    //                         width: "200px",
    //                         margin: "auto",
    //                         marginTop: "30pt"
    //                     }}
    //                     type="submit"
    //                 >
    //                     Login
    //                 </Button>
    //                 <p style={{ color: "var(--grey-500)", textAlign: "center", fontSize: "9pt", marginTop:"20pt" }}>
    //                     version: {process.env.REACT_APP_VERSION}
    //                 </p>
    //                 <p style={{ color: "var(--grey-500)", textAlign: "center", fontSize: "9pt" }}>
    //                     &copy;{new Date().getFullYear()} codebox
    //                 </p>
    //             </form>
    //         </div>
    //     </div>
    // )
    return (
        <React.Fragment>
            ddd
        </React.Fragment>
    )
}