# CMP Runbooks

## Emergency Certificate Rotation

### Scenario: Certificate compromised or expired

1. **Identify affected certificates**
   ```sql
   SELECT id, subject, not_after FROM certificates 
   WHERE subject LIKE '%affected-domain%' AND status = 'active';
   ```

2. **Revoke compromised certificates**
   ```bash
   curl -X POST http://api:8082/api/v1/certs/{cert_id}/revoke \
     -H "Authorization: Bearer $TOKEN" \
     -d '{"reason": "keyCompromise"}'
   ```

3. **Issue replacement certificates**
   - Use UI or API to request new certificates
   - Use different key material (different key size/algo)

4. **Deploy to agents**
   ```bash
   curl -X POST http://api:8082/api/v1/agents/{agent_id}/install \
     -H "Authorization: Bearer $TOKEN" \
     -d '{"cert_id": "new-cert-id", "path": "/path/to/cert.pem"}'
   ```

5. **Verify installation**
   - Check agent logs
   - Test HTTPS connections
   - Verify certificate chain

## Agent Offline Recovery

### Scenario: Agent not checking in or responding

1. **Check agent status**
   ```sql
   SELECT id, hostname, last_checkin, status FROM agents 
   WHERE id = 'agent-id';
   ```

2. **Verify network connectivity**
   - Ping agent host
   - Check firewall rules
   - Verify agent process is running

3. **Restart agent service**
   ```bash
   ssh agent-host
   systemctl restart cmp-agent
   systemctl status cmp-agent
   ```

4. **Force check-in**
   ```bash
   curl -X POST http://agent-host:8084/checkin
   ```

5. **Check agent logs**
   ```bash
   journalctl -u cmp-agent -f
   ```

6. **Re-register agent if needed**
   - Update database record
   - Issue new auth token from Vault
   - Update agent configuration

## Key Compromise Response

### Scenario: Private key exposed or suspected compromise

1. **Immediate Actions**
   - Revoke all certificates using compromised key
   - Mark key as compromised in system
   - Notify security team

2. **Certificate Revocation**
   ```bash
   # Revoke via API
   for cert_id in $(get-compromised-certs); do
     curl -X POST http://api:8082/api/v1/certs/$cert_id/revoke \
       -d '{"reason": "keyCompromise"}'
   done
   ```

3. **CA Notification**
   - If using external CA, notify them immediately
   - Request certificate revocation list update

4. **Audit Review**
   ```sql
   SELECT * FROM audit_logs 
   WHERE entity_type = 'certificate' 
   AND entity_id IN (compromised-ids)
   ORDER BY timestamp DESC;
   ```

5. **Key Rotation**
   - Generate new keys with higher security (4096-bit RSA or ECDSA)
   - Issue new certificates
   - Deploy immediately

6. **Post-Incident**
   - Root cause analysis
   - Update key generation policies
   - Review access controls

## Database Backup & Restore

### Backup
```bash
# Create backup
pg_dump -h db-host -U cmp_user -d cmp_db > backup-$(date +%Y%m%d).sql

# Compress
gzip backup-*.sql
```

### Restore
```bash
# Stop services
kubectl scale deployment cmp-issuer --replicas=0

# Restore database
gunzip < backup-*.sql.gz | psql -h db-host -U cmp_user -d cmp_db

# Verify data
psql -h db-host -U cmp_user -d cmp_db -c "SELECT COUNT(*) FROM certificates;"

# Restart services
kubectl scale deployment cmp-issuer --replicas=2
```

## Vault Unseal

### Scenario: Vault sealed after restart

1. **Check Vault status**
   ```bash
   vault status
   ```

2. **Unseal Vault** (requires unseal keys)
   ```bash
   vault operator unseal <key1>
   vault operator unseal <key2>
   vault operator unseal <key3>
   ```

3. **Verify PKI mount**
   ```bash
   vault secrets list
   vault read cmp-pki/config/urls
   ```

4. **Test certificate issuance**
   ```bash
   vault write cmp-pki/sign/cmp-role \
     common_name="test.example.com" \
     ttl="1h"
   ```

## High Availability Failover

### Scenario: Primary database failure

1. **Promote replica**
   ```sql
   -- On replica
   SELECT pg_promote();
   ```

2. **Update service configuration**
   - Update DB_HOST env var in deployments
   - Restart services to reconnect

3. **Verify connectivity**
   - Check service health endpoints
   - Verify API responses
   - Check audit logs

4. **Repair primary**
   - Fix primary database
   - Reconfigure as replica
   - Set up replication

## Certificate Renewal Automation

### Manual Renewal Process

1. **Identify expiring certificates**
   ```sql
   SELECT id, subject, not_after FROM certificates 
   WHERE not_after < NOW() + INTERVAL '30 days' 
   AND status = 'active';
   ```

2. **Request renewal**
   - Use same common name and SANs
   - Request via API or UI

3. **Verify renewal**
   - Check certificate not_after date
   - Test HTTPS connections
   - Verify certificate chain

### Automated Renewal (Planned)

- Job scheduler (Cron/K8s CronJob)
- Pre-renewal (30 days before expiry)
- Automatic installation
- Notification on failure

## Troubleshooting Common Issues

### Certificate Request Stuck in "pending"
1. Check adapter service logs
2. Verify Vault connectivity
3. Check adapter configuration in database
4. Review error messages in issuance_requests table

### Agent Installation Fails
1. Check agent connectivity to CMP API
2. Verify file path permissions
3. Check reload command syntax
4. Review agent logs

### High Latency
1. Check database query performance
2. Review Redis cache hit rate
3. Monitor service metrics
4. Check network latency between services
