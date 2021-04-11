import React from 'react';
import {Button} from "@material-ui/core";
import Request from "../request";

export default function PageAbout(props) {
    return (
        <div>
            About
            <Button onClick={() => {
                let req = {
                    "user_name": "admin",
                    "password": "admin"
                }
                Request('session_open', req)
                    .then((res) => {
                        if (res.status === 200) {
                            res.text().then(
                                (result) => {
                                    try {
                                        let obj = JSON.parse(result);
                                        console.log("session_open ok", obj)
                                    } catch (e) {
                                        console.log("session_open Wrong json", e)
                                    }

                                }
                            )
                        } else {
                            res.json().then(
                                (result) => {
                                    console.log("session_open ok", result)
                                }
                            );
                        }
                    });
            }}>Start</Button>
            <Button onClick={() => {
                let req = {
                }
                Request('public_channel_list', req)
                    .then((res) => {
                        if (res.status === 200) {
                            res.text().then(
                                (result) => {
                                    try {
                                        let obj = JSON.parse(result);
                                        console.log("session_open ok", obj)
                                    } catch (e) {
                                        console.log("session_open Wrong json", e)
                                    }

                                }
                            )
                        } else {
                            res.json().then(
                                (result) => {
                                    console.log("session_open ok", result)
                                }
                            );
                        }
                    });
            }}>Request</Button>
        </div>
    );
}
