import {Card, Flex, Input} from "antd";
import {useEffect, useState} from "react";
import {fetchWithAuth} from "../../utils/fetchWithAuth.js";

const CardContent = ({action, payload}) => {
    const decodedPayload = JSON.parse(atob(payload))
    return (
        <>
            <p>env: {decodedPayload.environment}</p>
            <p>project: {decodedPayload.project}</p>
            <p>release: {decodedPayload.release}</p>

            <p>---</p>

            {decodedPayload.items.map(row => (
                <p>{row.key}: {row.old_value} -> {row.new_value} </p>
            ))}
        </>
    )
}

const AuditList = () => {
    const [audits, setAudits] = useState([])

    useEffect(() => {
        fetchWithAuth('/api/v1/audits')
            .then((response) => {
                if (!response.ok) throw new Error("err");
                return response.json();
            })
            .then((data) => {
                setAudits(data.data.audits)
            })
    }, []);

    return (
        <>
            <Flex justify={"center"}>
                <Input placeholder="limit"/>
            </Flex>
            {audits.map(row =>(
                <Card
                    size="small"
                    title={row.action + " | " + (row.actor || "UNKNOWN") }
                    extra={row.ts}
                >
                    <CardContent action={row.action} payload={row.payload}/>
                </Card>
            ))}
        </>
    )
}

export default AuditList