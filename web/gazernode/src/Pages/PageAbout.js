import React, {useState} from 'react';
import Typography from "@material-ui/core/Typography";
import Request from "../request";

export default function PageAbout(props) {
    const [api ,setApi] = React.useState({})

    const requestApi = () => {
        let req = {
        }
        Request('service_api', req)
            .then((res) => {
                if (res.status === 200) {
                    res.json().then(
                        (result) => {
                            setApi(result)
                        }
                    );
                } else {
                    res.json().then(
                        (result) => {
                            //setErrorMessage(result.error)
                        }
                    );
                }
            });
    }

    const [firstRendering, setFirstRendering] = useState(true)
    if (firstRendering) {
        props.OnTitleUpdate("About")
        requestApi()
        setFirstRendering(false)
    }

    return (
        <div>
            <Typography style={{color: "#080", fontSize: "24pt"}}>{api.product}</Typography>
            <Typography style={{color: "#080", fontSize: "16pt"}}>version: {api.version}</Typography>
            <Typography style={{color: "#777", fontSize: "16pt"}}>{api.build_time}</Typography>
        </div>
    );
}
