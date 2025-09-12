import React, {useEffect, useMemo, useState, useCallback, useRef} from "react";
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
    DatePicker,
    Select,
    Descriptions
} from "antd";
import {SearchOutlined, ReloadOutlined} from "@ant-design/icons";
import {fetchWithAuth} from "../../utils/fetchWithAuth.js";
import dayjs from "dayjs";

const { RangePicker } = DatePicker;
const { Option } = Select;

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

// Компонент для изменений конфигурации
const ConfigUpdatedPayload = ({decoded}) => (
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
);

const DefaultPayload = ({decoded}) => {
    const entries = Object.entries(decoded || {}).filter(([key]) => key !== 'version');

    return (
        <Descriptions
            column={1}
            bordered
            size="small"
            labelStyle={{ fontWeight: 'bold', width: '120px' }}
        >
            {entries.map(([key, value]) => {
                let displayValue = value;
                if (typeof value === 'object' && value !== null) {
                    displayValue = JSON.stringify(value, null, 2);
                }

                return (
                    <Descriptions.Item key={key} label={key}>
                        {String(displayValue)}
                    </Descriptions.Item>
                );
            })}
        </Descriptions>
    );
};

const PrettyPayload = ({payload, action}) => {
    const decoded = useMemo(() => decodePayload(payload), [payload]);
    if (!decoded) return <div>Invalid payload</div>;

    switch (action) {
        case 'config_updated':
            return <ConfigUpdatedPayload decoded={decoded} />;
        default:
            return <DefaultPayload decoded={decoded} />;
    }
};

const Header = ({onRefresh, actor, setActor, action, setAction, dateRange, setDateRange, loading, actions}) => (
    <Row gutter={[8, 8]} align="middle" style={{marginBottom: 12}}>
        <Col xs={24} sm={12} md={8} lg={6}>
            <Input
                placeholder="Search by actor"
                value={actor}
                onChange={(e) => setActor(e.target.value)}
                prefix={<SearchOutlined />}
                allowClear
            />
        </Col>

        <Col xs={12} sm={6} md={4} lg={3}>
            <Select
                value={action}
                onChange={setAction}
                style={{width: '100%'}}
                loading={loading}
            >
                <Option value="">All Actions</Option>
                {actions.map(action => (
                    <Option key={action} value={action}>{action}</Option>
                ))}
            </Select>
        </Col>

        <Col xs={24} sm={12} md={8} lg={8}>
            <RangePicker
                value={dateRange}
                onChange={setDateRange}
                format="YYYY-MM-DD"
                style={{width: '100%'}}
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

    const actor = row.actor || "UNKNOWN";
    const ts = row.ts ? new Date(row.ts).toLocaleString() : "-";

    return (
        <Card
            size="small"
            style={{marginBottom: 12}}
            title={<span>{row.action} <span style={{color: 'var(--ant-gray-6)'}}> | {actor}</span></span>}
            extra={<small style={{color: 'var(--ant-gray-6)'}}>{ts}</small>}
        >
            <div style={{display: 'flex', gap: 12, flexDirection: 'column'}}>
                <PrettyPayload payload={payload} action={row.action} />
            </div>
        </Card>
    );
};

const AuditList = () => {
    const navigate = useNavigate();
    const [audits, setAudits] = useState([]);
    const [loading, setLoading] = useState(false);
    const [actionsLoading, setActionsLoading] = useState(false);
    const [actor, setActor] = useState('');
    const [action, setAction] = useState('');
    const [actions, setActions] = useState([]);
    const [dateRange, setDateRange] = useState([dayjs().subtract(1, 'day'), dayjs()]);
    const [error, setError] = useState(null);

    const fetchData = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const params = new URLSearchParams();
            if (action) params.set('action', action);
            if (actor) params.set('actor', actor);

            const [from, to] = dateRange;
            params.set('from', from.startOf('day').format('YYYY-MM-DDTHH:mm:ssZ'));
            params.set('to', to.endOf('day').format('YYYY-MM-DDTHH:mm:ssZ'));

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
    }, [actor, action, dateRange, navigate]);

    const fetchActions = useCallback(async () => {
        setActionsLoading(true);
        setError(null);
        try {
            const res = await fetchWithAuth("/api/v1/audits/actions", {}, navigate);
            if (!res.ok) throw new Error('fetch error');
            const json = await res.json();
            setActions(json.data || []);
        } catch (e) {
            setError('Failed to load actions');
            setActions([]);
        } finally {
            setActionsLoading(false);
        }
    }, [navigate]);

    const fetchDataRef = useRef(fetchData);
    fetchDataRef.current = fetchData;

    useEffect(() => {
        fetchDataRef.current();
        fetchActions();
    }, []);

    useEffect(() => {
        const t = setTimeout(() => {
            fetchDataRef.current();
        }, 400);
        return () => clearTimeout(t);
    }, [actor]);

    useEffect(() => {
        fetchDataRef.current();
    }, [action, dateRange]);

    return (
        <div>
            <Header
                onRefresh={fetchData}
                actor={actor}
                setActor={setActor}
                action={action}
                setAction={setAction}
                dateRange={dateRange}
                setDateRange={setDateRange}
                loading={loading || actionsLoading}
                actions={actions}
            />

            {(loading || actionsLoading) && <Spin style={{display: 'block', margin: '24px auto'}} />}

            {error && <div style={{color: 'var(--ant-error-color)', marginBottom: 12}}>{error}</div>}

            {!loading && !actionsLoading && audits.length === 0 && <Empty description="No audits" />}

            {audits.map((row, idx) => (
                <AuditCard key={row.id || idx} row={row} />
            ))}
        </div>
    );
};

export default AuditList;