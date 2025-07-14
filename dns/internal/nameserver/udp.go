package nameserver

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"strings"
	"sync"

	"golang.org/x/net/dns/dnsmessage"
)

// AuthoritativeServer
type AuthoritativeServer struct {
	logger *slog.Logger
	addr   *net.UDPAddr
	dht    dht
	wg     sync.WaitGroup
}

// NewAuthoritativeServer
func NewAuthoritativeServer(logger *slog.Logger, dht dht) (*AuthoritativeServer, error) {
	addr, err := net.ResolveUDPAddr("udp", "localhost:5300")
	if err != nil {
		return nil, fmt.Errorf("failed to resolve udp addr: %w", err)
	}
	return &AuthoritativeServer{
		addr:   addr,
		dht:    dht,
		logger: logger,
		wg:     sync.WaitGroup{},
	}, nil
}

// Listen
func (a *AuthoritativeServer) Listen(ctx context.Context) error {
	conn, err := net.ListenUDP("udp", a.addr)
	if err != nil {
		return fmt.Errorf("failed to create udp listener: %w", err)
	}
	defer func() {
		a.logger.Info("shutdown udp server")
		conn.Close()
		a.wg.Wait()
	}()

	a.logger.InfoContext(ctx, "authoritative server", slog.String("addr", a.addr.AddrPort().String()))

	go a.handleConn(ctx, conn)

	a.wg.Wait()
	<-ctx.Done()

	return nil
}

func (a *AuthoritativeServer) handleConn(ctx context.Context, conn *net.UDPConn) {
	a.wg.Add(1)
	defer a.wg.Done()

	for {
		var buf []byte = make([]byte, 512)

		_, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			a.logger.ErrorContext(ctx, "read from udp", slog.String("error", err.Error()))
			return
		}

		a.logger.InfoContext(
			ctx,
			"dns query",
			slog.String("zone", addr.Zone),
			slog.String("addr_port", addr.AddrPort().String()),
		)

		a.wg.Add(1)
		go func() {
			defer a.wg.Done()

			result, err := a.dnsLookup(buf)
			if err != nil {
				a.logger.ErrorContext(ctx, "dnslookup", slog.String("error", err.Error()))
				return
			}

			_, err = conn.WriteToUDP(result, addr)
			if err != nil {
				a.logger.ErrorContext(ctx, "write udp", slog.String("error", err.Error()))
			}
		}()
	}
}

func (a *AuthoritativeServer) dnsLookup(msg []byte) ([]byte, error) {
	var p = dnsmessage.Parser{}
	header, err := p.Start(msg)
	if err != nil {
	}

	qs, err := p.AllQuestions()
	if err != nil {
	}

	if len(qs) == 0 {

	}

	b := dnsmessage.NewBuilder(nil, header)
	b.EnableCompression()

	b.StartAnswers()
	err = a.answerQuestions(qs, b)
	if err != nil {

	}

	return b.Finish()
}

func (a *AuthoritativeServer) answerQuestions(qs []dnsmessage.Question, builder dnsmessage.Builder) error {
	for _, q := range qs {
		key := fmt.Sprintf("%s:%s", q.Name, q.Type)

		record, err := a.dht.getValue(key)
		if err == nil {
			switch strings.ToLower(record.RecordType) {
			case "a":
				err = addARecord(record, builder)
			case "cname":

			case "did":
			}

			if err != nil {
			}

			err = builder.Question(q)
			if err != nil {
			}
		}

		if err != nil && errors.Is(err, errKeyNotFound) {

		} else if err != nil {

		}
	}

	return nil
}

func addARecord(record *record, builder dnsmessage.Builder) error {
	return builder.AResource(dnsmessage.ResourceHeader{
		Name:  dnsmessage.MustNewName(record.Domain),
		Type:  dnsmessage.TypeA,
		Class: dnsmessage.ClassINET,
		TTL:   uint32(record.Ttl),
		// Length: ,
	}, dnsmessage.AResource{
		A: [4]byte(record.Value),
	})
}
