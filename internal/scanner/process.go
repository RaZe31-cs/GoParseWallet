package scanner

import (
	"ton-lessons2/internal/app"

	"github.com/sirupsen/logrus"
	"github.com/xssnick/tonutils-go/tlb"
)



func processJetton(
	internalMessage *tlb.InternalMessage,
) error {
	bodySlice := internalMessage.Body.BeginParse()

	if internalMessage.Body == nil {
		logrus.Warn("empty body")
		return nil
	}

	opcode, err := bodySlice.LoadUInt(32)
	if err != nil {
		logrus.Error("error when get opcode: ", err)
		return err

	}

	if opcode != 0x7362d09c {
		logrus.Warn("not jetton notification")
		return err

	}

	queryId, err := bodySlice.LoadUInt(64)
	if err != nil {
		logrus.Error("query id err: ", err)
		return err

	}

	amount, err := bodySlice.LoadCoins()
	if err != nil {
		logrus.Error("amount err: ", err)
		return err

	}

	sender, err := bodySlice.LoadAddr()
	if err != nil {
		logrus.Error("address err: ", err)
		return err

	}

	fwdPayload, err := bodySlice.LoadMaybeRef()
	if err != nil {
		logrus.Error("fwd payload err: ", err)
		return err
	}

	fwdOp, err := fwdPayload.LoadUInt(32)
	if err != nil {
		logrus.Error("fwd op err: ", err)
		return err
	}

	if fwdOp != 0 {
		logrus.Error("not text comment")
		return err
	}

	comment, err := fwdPayload.LoadStringSnake()
	if err != nil {
		logrus.Error("text comment err: ", err)
		return err
	}

	if comment != app.CFG.Wallet.UuidGoogle {
		logrus.Warn("[JTN] Comment != uuid")
	} else {
		logrus.Info("[JTN] new transaction!")
		logrus.Info("[JTN] sender: ", sender)
		logrus.Info("[JTN] amount: ", amount)
		logrus.Info("[JTN] query id: ", queryId)
		logrus.Info("[JTN] comment: ", comment)
	}
	return nil
}


func processTon(
	internalMessage *tlb.InternalMessage,
) error {
	comment := internalMessage.Comment()
	logrus.Info(comment)
	amount := internalMessage.Amount
	sender := internalMessage.SenderAddr()
	if comment != app.CFG.Wallet.UuidGoogle {
		logrus.Warn("[TON] Comment != uuid")
	} else {
		logrus.Info("[TON] new transaction!")
		logrus.Info("[TON] Comment: ", comment)
		logrus.Info("[TON] Amount: ", amount)
		logrus.Info("[TON] Sender: ", sender)
	}
	return nil
}


func (s *scanner) processTransaction(
	trans *tlb.Transaction,
) error {
	if trans.IO.In.MsgType != tlb.MsgTypeInternal {
		return nil
	}

	internalMessage := trans.IO.In.AsInternal()
	if internalMessage.DestAddr().String() != app.CFG.Wallet.AddressParse {
		return nil
	}

	logrus.Info("SrcAddress: ", internalMessage.SrcAddr.String())
	err := processJetton(internalMessage)
	if err != nil {
		err = processTon(internalMessage)
		if err != nil {
			return err
		}
	}
	return nil
}
