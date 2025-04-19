import React, { useCallback, useEffect, useState } from "react";
import { Container, Input, InputGroup, InputGroupText } from "reactstrap";
import CodeboxLogo from "../../assets/images/codebox-logo-white.png";
import { Http } from "../../api/http";
import { RequestStatus } from "../../api/types";
import { toast, ToastContainer } from "react-toastify";

// TODO: improve login process use this token as
// temporarily token to request a real token
export default function CliLogin() {
  const [token, setToken] = useState<string>("");

  const requestToken = useCallback(async () => {
    var [status, statusCode, data] = await Http.Request(
      `${Http.GetServerURL()}/api/v1/auth/cli-login`,
      "POST",
      null
    );
    if (status === RequestStatus.OK && statusCode === 200) {
      setToken(data.token);
    } else {
      toast.error(`Cannot request new token, ${statusCode}`);
    }
  }, []);

  useEffect(() => {
    requestToken();
  }, [requestToken]);

  return (
    <React.Fragment>
      {token.length > 0 && (
        <div className="page page-center">
          <Container className="w-100 d-flex justify-content-center">
            <div className="mb-5">
              <div className="section-header">
                <span className="d-flex align-items-center w-100 justify-content-center">
                  <img src={CodeboxLogo} alt="logo" width={185} />
                </span>
                <p className="section-description text-secondary text-center mt-4">
                  Copy this token, then return to the <br />
                  CLI/extension and paste it.
                </p>
                <div className="mt-5 btn-list d-flex justify-content-center">
                  <InputGroup>
                    <Input type={"password"} value={token} disabled />
                    <InputGroupText
                      style={{ cursor: "pointer" }}
                      onClick={async() => {
                        await navigator.clipboard.writeText(token);
                        toast.info("Token has been copied to clipboard");
                      }}
                    >
                      Copy
                    </InputGroupText>
                  </InputGroup>
                </div>
              </div>
            </div>
          </Container>
        </div>
      )}
      <ToastContainer/>
    </React.Fragment>
  );
}
