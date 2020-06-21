package plug

import (
	"fmt"
	"log"

	data "UwdBot/database"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	usersInTopList = 10
	robCount       = 2
)

type CaravanRobber struct {
	UserID   uint64
	Username string
	Power    int
}

type CaravanRobbers [robCount]CaravanRobber

func (c *CaravanRobbers) checkRobbersCount() int {
	count := 0
	for _, caravan := range c {
		if caravan.UserID != 0 {
			count++
		}
	}
	return count
}

func (w *Wars) RobCaravans(msg *tgbotapi.Message, user *data.User) string {
	robbersCount := w.robbers.checkRobbersCount()
	if robbersCount == robCount {
		return "🐫🐪🐫"
	}

	w.robbers[robbersCount] = CaravanRobber{
		user.UserID, user.Username, user.WeaponsPower,
	}
	robbersCount = w.robbers.checkRobbersCount()
	if robbersCount == robCount {
		return w.caravansStart()
	}

	return fmt.Sprintf(
		"Для отправления каравана нужно еще ***%d*** грабителя!",
		robCount-robbersCount,
	)

}

func (w *Wars) caravansStart() string {
	return "Начинаем набег на караван!"
}

func (w *Wars) GetTopPlayers(count int) string {
	user := data.User{}
	result := "**ТОП ИГРОКОВ:**\n"
	topUsers, err := user.GetTopUsers(count)

	log.Println(err)

	if err != nil {
		return "Что то пошло не так..."
	}

	for i, us := range topUsers {
		result += fmt.Sprintf(
			"%d) %s: %d👑 %d💰\n",
			i+1,
			GetMarkdownUsername(us.Username),
			us.Reputation,
			us.Coins,
		)
	}

	result += "\n__Региструйся и победи всех__ **/reg**"
	return result
}

func (w *Wars) GetShop(msg *tgbotapi.Message) string {
	weap := data.Weapon{}
	weapons, err := weap.GetAllWeapons()
	if err != nil {
		return "Не удалось загрузить магазин..."
	}
	reply := "***Уютный shop 🛒 ***\n\n***Оружие:***\n"
	for _, w := range weapons {
		reply += fmt.Sprintf(
			"%d) ___%s___ %d🗡️, %d💰\n",
			w.ID,
			w.Name,
			w.Power,
			w.Cost,
		)
	}
	reply += "\n___Интересный стафф:___\nПоявится в скором времени...\n\n___Купить товар - реплай на сообщение buy номер товара___"
	return reply
}

func (w *Wars) buyItem(item int, msg *tgbotapi.Message) {
	var err error
	var user data.User
	user, err = user.FindUserByID(msg.From.ID)
	if err != nil {
		w.c.SendReplyToMessage(msg, "Вы не зарегистрированы /reg")
		return
	}
	var weapon data.Weapon
	weapon, err = weapon.GetWeaponsByID(item)
	if err != nil {
		w.c.SendReplyToMessage(msg, "Некоректный номер товара!")
		return
	}

	if user.Coins >= weapon.Cost {
		user.DecreaseMoney(weapon.Cost)
		user.AddPower(weapon.Power)
		w.c.SendMarkdownReply(
			msg,
			fmt.Sprintf(
				"Списано ***%d***💰, куплен(а): ___%s___!\n\n Прибавлено %d к боевой мощи!",
				weapon.Cost,
				weapon.Name,
				weapon.Power,
			),
		)
	} else {
		w.c.SendMarkdownReply(
			msg,
			fmt.Sprintf(
				"Вам не хватает ***%d***💰, чтобы купить ___%s___!",
				weapon.Cost-user.Coins,
				weapon.Name,
			),
		)
	}
}
