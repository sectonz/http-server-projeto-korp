# HTTP Server - Projeto Korp

Projeto desenvolvido como solução para o desafio técnico de DevOps/SRE e Backend. A solução contempla a criação de um serviço HTTP em Golang, containerização com Docker, orquestração com Docker Compose, proxy reverso com Nginx, observabilidade completa com Prometheus e Grafana (com provisionamento automatizado) e automação de infraestrutura via Ansible.

---

## Visão Geral da Solução

A arquitetura do projeto foi estruturada em três camadas principais:

1. **Aplicação (Go)**: Servidor HTTP (`http-server-projeto-korp`) com endpoints de negócio (`/projeto-korp`), healthcheck (`/health`) e métricas padrão Prometheus (`/metrics`).
2. **Infraestrutura e Proxy (Docker & Nginx)**: Roteador reverso Nginx na porta `80` repassando o tráfego para a aplicação isolada na rede bridge `projeto-korp-network`.
3. **Observabilidade (Prometheus & Grafana)**: Coleta automática de métricas de disponibilidade, volume e latência pelo Prometheus e exibição em tempo real através de um dashboard provisionado como código no Grafana.
4. **Automação (Ansible)**: Automação ponta a ponta para instalação do Docker, build das imagens, execução da stack via Compose e validação do serviço.

---

## Arquitetura do Ambiente

```text
               +-------------------------------------------------------------+
               |                    Docker Host / Client                     |
               +-------------------------------------------------------------+
                                              │
                      ┌───────────────────────┼───────────────────────┐
                      │ (Porta 80)            │ (Porta 9090)          │ (Porta 3000)
                      ▼                       ▼                       ▼
             ┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
             │      Nginx      │    │   Prometheus    │    │     Grafana     │
             │ (Reverse Proxy) │    │  (Scrape / 5s)  │    │  (Dashboards)   │
             └────────┬────────┘    └────────┬────────┘    └────────┬────────┘
                      │                      │                      │
                      │                      │ Consulta             │ Consome
                      │                      ▼                      │
                      │             ┌─────────────────┐             │
                      └────────────►│   Go Server     │◄────────────┘
                                    │ (Porta 8080)    │
                                    └─────────────────┘
                                       Rede Bridge:
                                  projeto-korp-network
```

---

## Tecnologias Utilizadas

- **Linguagem**: [Go 1.23](https://golang.org/)
- **Containerização**: [Docker](https://www.docker.com/)
- **Web Server / Proxy Reverso**: [Nginx](https://hub.docker.com/_/nginx)
- **Coleta de Métricas**: [Prometheus](https://prometheus.io/)
- **Visualização & Dashboards**: [Grafana](https://grafana.com/)
- **Automação de Infraestrutura**: [Ansible](https://www.ansible.com/)
