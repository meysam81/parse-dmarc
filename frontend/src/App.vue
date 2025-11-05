<template>
  <div class="app">
    <header class="header">
      <div class="container">
        <h1 class="title">
          <span class="icon">🛡️</span>
          DMARC Report Dashboard
        </h1>
        <p class="subtitle">Email Authentication & Compliance Monitoring</p>
      </div>
    </header>

    <main class="main">
      <div class="container">
        <!-- Statistics Cards -->
        <div class="stats-grid" v-if="statistics">
          <div class="stat-card">
            <div class="stat-icon">📊</div>
            <div class="stat-content">
              <div class="stat-value">{{ statistics.total_reports }}</div>
              <div class="stat-label">Total Reports</div>
            </div>
          </div>

          <div class="stat-card">
            <div class="stat-icon">📧</div>
            <div class="stat-content">
              <div class="stat-value">{{ formatNumber(statistics.total_messages) }}</div>
              <div class="stat-label">Total Messages</div>
            </div>
          </div>

          <div class="stat-card">
            <div class="stat-icon">✅</div>
            <div class="stat-content">
              <div class="stat-value">{{ statistics.compliance_rate.toFixed(1) }}%</div>
              <div class="stat-label">Compliance Rate</div>
            </div>
          </div>

          <div class="stat-card">
            <div class="stat-icon">🌐</div>
            <div class="stat-content">
              <div class="stat-value">{{ statistics.unique_source_ips }}</div>
              <div class="stat-label">Unique Sources</div>
            </div>
          </div>
        </div>

        <!-- Top Sources -->
        <div class="section">
          <h2 class="section-title">Top Sending Sources</h2>
          <div class="card">
            <div class="source-list" v-if="topSources.length > 0">
              <div
                v-for="source in topSources"
                :key="source.source_ip"
                class="source-item"
              >
                <div class="source-ip">{{ source.source_ip }}</div>
                <div class="source-stats">
                  <div class="source-count">{{ formatNumber(source.count) }} messages</div>
                  <div class="source-bar">
                    <div
                      class="source-bar-pass"
                      :style="{ width: getPassPercentage(source) + '%' }"
                    ></div>
                    <div
                      class="source-bar-fail"
                      :style="{ width: getFailPercentage(source) + '%' }"
                    ></div>
                  </div>
                  <div class="source-legend">
                    <span class="legend-pass">{{ source.pass }} pass</span>
                    <span class="legend-fail">{{ source.fail }} fail</span>
                  </div>
                </div>
              </div>
            </div>
            <div v-else class="empty-state">
              No data available yet. Reports will appear here once fetched.
            </div>
          </div>
        </div>

        <!-- Recent Reports -->
        <div class="section">
          <h2 class="section-title">Recent Reports</h2>
          <div class="card">
            <div class="table-container" v-if="reports.length > 0">
              <table class="report-table">
                <thead>
                  <tr>
                    <th>Organization</th>
                    <th>Domain</th>
                    <th>Date Range</th>
                    <th>Messages</th>
                    <th>Compliance</th>
                    <th>Policy</th>
                  </tr>
                </thead>
                <tbody>
                  <tr
                    v-for="report in reports"
                    :key="report.id"
                    class="report-row"
                    @click="viewReport(report)"
                  >
                    <td>{{ report.org_name }}</td>
                    <td><code>{{ report.domain }}</code></td>
                    <td class="date-cell">{{ formatDate(report.date_begin) }}</td>
                    <td>{{ formatNumber(report.total_messages) }}</td>
                    <td>
                      <span
                        class="compliance-badge"
                        :class="getComplianceClass(report.compliance_rate)"
                      >
                        {{ report.compliance_rate.toFixed(1) }}%
                      </span>
                    </td>
                    <td>
                      <span class="policy-badge">{{ report.policy_p }}</span>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
            <div v-else class="empty-state">
              No reports available yet. Check your IMAP configuration and run the fetch process.
            </div>
          </div>
        </div>

        <!-- Report Detail Modal -->
        <div v-if="selectedReport" class="modal" @click="closeModal">
          <div class="modal-content" @click.stop>
            <div class="modal-header">
              <h3>Report Details</h3>
              <button class="modal-close" @click="closeModal">×</button>
            </div>
            <div class="modal-body">
              <div class="detail-grid">
                <div class="detail-item">
                  <strong>Organization:</strong> {{ selectedReport.report_metadata.org_name }}
                </div>
                <div class="detail-item">
                  <strong>Domain:</strong> {{ selectedReport.policy_published.domain }}
                </div>
                <div class="detail-item">
                  <strong>Report ID:</strong> {{ selectedReport.report_metadata.report_id }}
                </div>
                <div class="detail-item">
                  <strong>Policy:</strong> {{ selectedReport.policy_published.p }}
                </div>
              </div>

              <h4 class="detail-subtitle">Records ({{ selectedReport.records.length }})</h4>
              <div class="records-list">
                <div
                  v-for="(record, idx) in selectedReport.records"
                  :key="idx"
                  class="record-item"
                >
                  <div class="record-header">
                    <span class="record-ip">{{ record.row.source_ip }}</span>
                    <span class="record-count">{{ record.row.count }} messages</span>
                  </div>
                  <div class="record-details">
                    <span :class="'result-badge ' + (record.row.policy_evaluated.dkim === 'pass' ? 'pass' : 'fail')">
                      DKIM: {{ record.row.policy_evaluated.dkim }}
                    </span>
                    <span :class="'result-badge ' + (record.row.policy_evaluated.spf === 'pass' ? 'pass' : 'fail')">
                      SPF: {{ record.row.policy_evaluated.spf }}
                    </span>
                    <span class="result-badge">
                      Disposition: {{ record.row.policy_evaluated.disposition }}
                    </span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>

    <footer class="footer">
      <div class="container">
        <p>Built with Vue 3 + Go | RFC 7489 Compliant</p>
      </div>
    </footer>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'

export default {
  name: 'App',
  setup() {
    const statistics = ref(null)
    const topSources = ref([])
    const reports = ref([])
    const selectedReport = ref(null)
    const loading = ref(true)

    const fetchStatistics = async () => {
      try {
        const response = await fetch('/api/statistics')
        statistics.value = await response.json()
      } catch (error) {
        console.error('Failed to fetch statistics:', error)
      }
    }

    const fetchTopSources = async () => {
      try {
        const response = await fetch('/api/top-sources?limit=10')
        topSources.value = await response.json()
      } catch (error) {
        console.error('Failed to fetch top sources:', error)
      }
    }

    const fetchReports = async () => {
      try {
        const response = await fetch('/api/reports?limit=20')
        reports.value = await response.json()
      } catch (error) {
        console.error('Failed to fetch reports:', error)
      }
    }

    const viewReport = async (report) => {
      try {
        const response = await fetch(`/api/reports/${report.id}`)
        selectedReport.value = await response.json()
      } catch (error) {
        console.error('Failed to fetch report details:', error)
      }
    }

    const closeModal = () => {
      selectedReport.value = null
    }

    const formatNumber = (num) => {
      return new Intl.NumberFormat().format(num)
    }

    const formatDate = (timestamp) => {
      return new Date(timestamp * 1000).toLocaleDateString()
    }

    const getPassPercentage = (source) => {
      const total = source.count
      return total > 0 ? (source.pass / total * 100) : 0
    }

    const getFailPercentage = (source) => {
      const total = source.count
      return total > 0 ? (source.fail / total * 100) : 0
    }

    const getComplianceClass = (rate) => {
      if (rate >= 95) return 'high'
      if (rate >= 70) return 'medium'
      return 'low'
    }

    const loadData = async () => {
      loading.value = true
      await Promise.all([
        fetchStatistics(),
        fetchTopSources(),
        fetchReports()
      ])
      loading.value = false
    }

    onMounted(() => {
      loadData()
      // Auto-refresh every 5 minutes
      setInterval(loadData, 5 * 60 * 1000)
    })

    return {
      statistics,
      topSources,
      reports,
      selectedReport,
      loading,
      viewReport,
      closeModal,
      formatNumber,
      formatDate,
      getPassPercentage,
      getFailPercentage,
      getComplianceClass
    }
  }
}
</script>

<style scoped>
.app {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.header {
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
  padding: 2rem 0;
  color: white;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 2rem;
}

.title {
  font-size: 2.5rem;
  font-weight: 700;
  margin-bottom: 0.5rem;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.icon {
  font-size: 2rem;
}

.subtitle {
  font-size: 1.1rem;
  opacity: 0.9;
}

.main {
  flex: 1;
  padding: 2rem 0;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.stat-card {
  background: white;
  border-radius: 12px;
  padding: 1.5rem;
  display: flex;
  align-items: center;
  gap: 1rem;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  transition: transform 0.2s;
}

.stat-card:hover {
  transform: translateY(-4px);
}

.stat-icon {
  font-size: 2.5rem;
}

.stat-content {
  flex: 1;
}

.stat-value {
  font-size: 2rem;
  font-weight: 700;
  color: #667eea;
}

.stat-label {
  color: #666;
  font-size: 0.9rem;
}

.section {
  margin-bottom: 2rem;
}

.section-title {
  color: white;
  font-size: 1.5rem;
  margin-bottom: 1rem;
  font-weight: 600;
}

.card {
  background: white;
  border-radius: 12px;
  padding: 1.5rem;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.source-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.source-item {
  padding: 1rem;
  background: #f8f9fa;
  border-radius: 8px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 1rem;
}

.source-ip {
  font-family: 'Courier New', monospace;
  font-weight: 600;
  color: #333;
  min-width: 150px;
}

.source-stats {
  flex: 1;
}

.source-count {
  font-size: 0.9rem;
  color: #666;
  margin-bottom: 0.5rem;
}

.source-bar {
  height: 8px;
  background: #e0e0e0;
  border-radius: 4px;
  overflow: hidden;
  display: flex;
  margin-bottom: 0.5rem;
}

.source-bar-pass {
  background: #4caf50;
  transition: width 0.3s;
}

.source-bar-fail {
  background: #f44336;
  transition: width 0.3s;
}

.source-legend {
  display: flex;
  gap: 1rem;
  font-size: 0.85rem;
}

.legend-pass {
  color: #4caf50;
}

.legend-fail {
  color: #f44336;
}

.table-container {
  overflow-x: auto;
}

.report-table {
  width: 100%;
  border-collapse: collapse;
}

.report-table th {
  text-align: left;
  padding: 0.75rem;
  background: #f8f9fa;
  font-weight: 600;
  color: #333;
  border-bottom: 2px solid #dee2e6;
}

.report-table td {
  padding: 0.75rem;
  border-bottom: 1px solid #dee2e6;
}

.report-row {
  cursor: pointer;
  transition: background 0.2s;
}

.report-row:hover {
  background: #f8f9fa;
}

.date-cell {
  color: #666;
  font-size: 0.9rem;
}

code {
  background: #f8f9fa;
  padding: 0.2rem 0.4rem;
  border-radius: 4px;
  font-family: 'Courier New', monospace;
  font-size: 0.9rem;
}

.compliance-badge {
  padding: 0.25rem 0.75rem;
  border-radius: 12px;
  font-size: 0.85rem;
  font-weight: 600;
}

.compliance-badge.high {
  background: #d4edda;
  color: #155724;
}

.compliance-badge.medium {
  background: #fff3cd;
  color: #856404;
}

.compliance-badge.low {
  background: #f8d7da;
  color: #721c24;
}

.policy-badge {
  padding: 0.25rem 0.75rem;
  border-radius: 12px;
  background: #e7f3ff;
  color: #0056b3;
  font-size: 0.85rem;
  font-weight: 600;
  text-transform: uppercase;
}

.empty-state {
  text-align: center;
  padding: 3rem;
  color: #666;
}

.modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 2rem;
}

.modal-content {
  background: white;
  border-radius: 12px;
  max-width: 800px;
  width: 100%;
  max-height: 80vh;
  overflow-y: auto;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.modal-header {
  padding: 1.5rem;
  border-bottom: 1px solid #dee2e6;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.modal-header h3 {
  margin: 0;
  font-size: 1.5rem;
}

.modal-close {
  background: none;
  border: none;
  font-size: 2rem;
  cursor: pointer;
  color: #999;
  padding: 0;
  width: 2rem;
  height: 2rem;
  display: flex;
  align-items: center;
  justify-content: center;
}

.modal-close:hover {
  color: #333;
}

.modal-body {
  padding: 1.5rem;
}

.detail-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
  margin-bottom: 1.5rem;
}

.detail-item {
  padding: 0.75rem;
  background: #f8f9fa;
  border-radius: 6px;
}

.detail-subtitle {
  margin: 1.5rem 0 1rem;
  font-size: 1.2rem;
}

.records-list {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.record-item {
  border: 1px solid #dee2e6;
  border-radius: 6px;
  padding: 1rem;
}

.record-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 0.75rem;
  font-weight: 600;
}

.record-ip {
  font-family: 'Courier New', monospace;
  color: #667eea;
}

.record-count {
  color: #666;
  font-size: 0.9rem;
}

.record-details {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}

.result-badge {
  padding: 0.25rem 0.75rem;
  border-radius: 12px;
  font-size: 0.85rem;
  background: #e9ecef;
  color: #495057;
}

.result-badge.pass {
  background: #d4edda;
  color: #155724;
}

.result-badge.fail {
  background: #f8d7da;
  color: #721c24;
}

.footer {
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
  padding: 1.5rem 0;
  color: white;
  text-align: center;
  margin-top: auto;
}

@media (max-width: 768px) {
  .title {
    font-size: 1.75rem;
  }

  .stats-grid {
    grid-template-columns: 1fr;
  }

  .source-item {
    flex-direction: column;
    align-items: flex-start;
  }

  .source-ip {
    min-width: auto;
  }
}
</style>
