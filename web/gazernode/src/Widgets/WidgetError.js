import React from "react";
import Grid from "@material-ui/core/Grid";

const formatTime = (date) => {

    let hh = date.getHours();
    if (hh < 10) hh = '0' + hh;

    let mm = date.getMinutes();
    if (mm < 10) mm = '0' + mm;

    let ss = date.getSeconds();
    if (ss < 10) ss = '0' + ss;

    return hh + ':' + mm + ':' + ss;
}

export default function WidgetError(props) {
    let timeString = formatTime(new Date())
    if (props.Message === undefined || props.Message === "") {
        return (
            <Grid container direction="row" style={{padding: "10px"}}>
                <Grid item>
                </Grid>
                <Grid item>
                    <Grid container direction="column" style={{padding: "10px"}}>
                        <Grid item>
                            <div style={{color: "#080", fontSize: "12pt"}}>connected</div>
                        </Grid>
                        <Grid item>
                            <div style={{color: "#444", fontSize: "10pt"}}>{timeString}</div>
                        </Grid>
                    </Grid>
                </Grid>
            </Grid>
        );
    }

    return (
        <Grid container direction="row" style={{padding: "10px"}}>
            <Grid item>
                <div><img src="https://gazer.cloud/channel_files/waiting.svg" alt={"waiting"}/></div>
            </Grid>
            <Grid item>
                <Grid container direction="column" style={{padding: "10px"}}>
                    <Grid item>
                        <div style={{color: "#F60", fontSize: "12pt"}}>{props.Message}</div>
                    </Grid>
                    <Grid item>
                        <div style={{color: "#444", fontSize: "10pt"}}>{timeString}</div>
                    </Grid>
                </Grid>
            </Grid>
        </Grid>
    );
}
