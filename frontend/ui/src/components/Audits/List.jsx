import React, {useEffect, useMemo, useState} from "react";
import { useNavigate } from "react-router-dom";
import {
    Card,
    Input,
    Tag,
    Space,
    Empty,
    Spin,
    Button,
    Row,
    Col,
} from "antd";
import {SearchOutlined, ReloadOutlined} from "@ant-design/icons";
import {fetchWithAuth} from "../../utils/fetchWithAuth.js";

const decodePayload = (payload) => {
    try {
        if (!payload) return null;
        const decoded = atob(payload);
        return JSON.parse(decoded);
    } catch (e) {
        return null;
    }
};

const DiffLine = ({label, oldValue, newValue}) => {
    const oldStr = oldValue === null || oldValue === undefined ? "" : String(oldValue);
    const newStr = newValue === null || newValue === undefined ? "" : String(newValue);
    const changed = oldStr !== newStr;

    return (
        <div style={{display: 'flex', gap: 16, alignItems: 'flex-start'}}>
            <div style={{minWidth: 220, whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis'}}>
                <b>{label}</b>
            </div>

            <div style={{flex: 1}}>
                <div style={{display: 'flex', flexDirection: 'column'}}>
                    <div style={{marginBottom: 14}}>
                        <div style={{marginTop: 6, padding: 8, borderRadius: 6, background: changed ? 'rgba(255,230,230,0.6)' : 'transparent'}}>
                            <code style={{display: 'block', overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'pre-wrap', textDecoration: changed ? 'line-through' : 'none', color: changed ? '#cf1322' : 'inherit'}}>{oldStr}</code>
                        </div>
                    </div>

                    <div style={{height: 8}} />

                    <div>
                        <div style={{marginTop: 6, padding: 8, borderRadius: 6, background: changed ? 'rgba(230,255,230,0.6)' : 'transparent'}}>
                            <code style={{display: 'block', overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'pre-wrap', fontWeight: changed ? 700 : 400, color: changed ? '#237804' : 'inherit'}}>{newStr}</code>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
};

const PrettyPayload = ({payload}) => {
    const decoded = useMemo(() => decodePayload(payload), [payload]);
    if (!decoded) return <div>Invalid payload</div>;

    return (
        <div className="audit-payload">
            <Space direction="vertical" size={12} style={{width: "100%"}}>
                <div>
                    <Tag>env: {decoded.environment || "-"}</Tag>
                    <Tag>project: {decoded.project || "-"}</Tag>
                    <Tag>release: {decoded.release || "-"}</Tag>
                </div>

                <div style={{display: "grid", gap: 12}}>
                    {(decoded.items || []).map((it, i) => (
                        <DiffLine key={i} label={it.key} oldValue={it.old_value} newValue={it.new_value} />
                    ))}
                </div>
            </Space>
        </div>
    );
};

const Header = ({onRefresh, limit, setLimit, query, setQuery, loading}) => (
    <Row gutter={[8, 8]} align="middle" style={{marginBottom: 12}}>
        <Col xs={24} sm={12} md={8} lg={6}>
            <Input
                placeholder="Search by actor / action / project"
                value={query}
                onChange={(e) => setQuery(e.target.value)}
                prefix={<SearchOutlined />}
                allowClear
            />
        </Col>

        <Col xs={12} sm={6} md={4} lg={3}>
            <Input
                placeholder="limit"
                value={limit}
                onChange={(e) => setLimit(e.target.value.replace(/[^0-9]/g, ""))}
            />
        </Col>

        <Col>
            <Space>
                <Button icon={<ReloadOutlined />} onClick={onRefresh} loading={loading} />
            </Space>
        </Col>
    </Row>
);

const AuditCard = ({row}) => {
    const payload = row.payload;
    const decoded = decodePayload(payload);

    const actor = row.actor || "UNKNOWN";
    const ts = row.ts ? new Date(row.ts).toLocaleString() : "-";

    return (
        <Card
            size="small"
            style={{marginBottom: 12}}
            title={<span>{row.action} <span style={{color: 'var(--ant-gray-6)', marginLeft: 8}}> | {actor}</span></span>}
            extra={<small style={{color: 'var(--ant-gray-6)'}}>{ts}</small>}
        >
            <div style={{display: 'flex', gap: 12, flexDirection: 'column'}}>
                <PrettyPayload payload={payload} />
            </div>
        </Card>
    );
};

const AuditList = () => {
    const navigate = useNavigate();
    const [audits, setAudits] = useState([]);
    const [loading, setLoading] = useState(false);
    const [limit, setLimit] = useState('50');
    const [query, setQuery] = useState('');
    const [error, setError] = useState(null);

    const fetchData = async () => {
        setLoading(true);
        setError(null);
        try {
            const params = new URLSearchParams();
            if (limit) params.set('limit', String(Number(limit) || 50));
            if (query) params.set('q', query);

            const res = await fetchWithAuth(`/api/v1/audits?${params.toString()}`, {}, navigate);
            if (!res.ok) throw new Error('fetch error');
            const json = await res.json();
            setAudits(json.data?.audits || []);
        } catch (e) {
            setError('Failed to load data');
            setAudits([]);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchData();
    }, []);

    useEffect(() => {
        const t = setTimeout(() => {
            fetchData();
        }, 400);
        return () => clearTimeout(t);
    }, [limit, query]);

    return (
        <div>
            <Header onRefresh={fetchData} limit={limit} setLimit={setLimit} query={query} setQuery={setQuery} loading={loading} />

            {loading && <Spin style={{display: 'block', margin: '24px auto'}} />}

            {error && <div style={{color: 'var(--ant-error-color)', marginBottom: 12}}>{error}</div>}

            {!loading && audits.length === 0 && <Empty description="No audits" />}

            {audits.map((row, idx) => (
                <AuditCard key={row.id || idx} row={row} />
            ))}
        </div>
    );
};

export default AuditList;
