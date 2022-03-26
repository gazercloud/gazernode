import React, {useState} from 'react';
import {Button} from "@material-ui/core";
import Request from "../request";

function getCookie(name) {
    let matches = document.cookie.match(new RegExp(
        "(?:^|; )" + name.replace(/([\.$?*|{}\(\)\[\]\\\/\+^])/g, '\\$1') + "=([^;]*)"
    ));
    return matches ? decodeURIComponent(matches[1]) : undefined;
}

export default function PageAccount(props) {

    const [firstRendering, setFirstRendering] = useState(true)
    if (firstRendering) {
        props.OnTitleUpdate("Account")
        setFirstRendering(false)
    }

    return (
        <div>
            <Button onClick={() => {
                let req = {
                    session_token: getCookie('session_token')
                }
                Request('session_remove', req)
                    .then((res) => {
                        if (res.status === 200) {
                            props.OnNeedUpdate()
                        } else {
                        }
                    });
            }}>LogOut</Button>
        </div>
    );
}
