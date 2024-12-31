package user

import (
	"context"
	"crypto/rand"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"

	"github.com/mpu-cad/gw-backend-go/internal/entity"
	"github.com/mpu-cad/gw-backend-go/internal/logger"
	"github.com/mpu-cad/gw-backend-go/internal/models"
)

type UCUser struct {
	user   userRepos
	mailer mailer
	redis  redis
}

func NewUCUser(user userRepos, mailer mailer, redis redis) *UCUser {
	return &UCUser{
		user:   user,
		mailer: mailer,
		redis:  redis,
	}
}

func (u *UCUser) Registration(ctx context.Context, request models.User) (*int, error) {
	pass, err := u.createHashPassword(request.HashPass)
	if err != nil {
		return nil, err
	}

	request.HashPass = pass

	id, err := u.user.InsertUser(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("can not insert user, err: %w", err)
	}

	if id == nil {
		return nil, fmt.Errorf("can not insert user, err: %w", err)
	}

	code := generateVerificationCode()
	u.redis.SaveUsersRegistrationCode(ctx, code, *id)

	go func(code string) {
		content := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Подтверждение регистрации</title>
    <style>
        @media only screen and (max-width: 620px) {
            table[class="body"] h1 {
                font-size: 28px !important;
                margin-bottom: 10px !important;
            }

            table[class="body"] .wrapper,
            table[class="body"] .article {
                padding: 10px !important;
            }

            table[class="body"] .content {
                padding: 0 !important;
            }

            table[class="body"] .container {
                padding: 0 !important;
                width: 100% !important;
            }

            table[class="body"] .main {
                border-left-width: 0 !important;
                border-radius: 0 !important;
                border-right-width: 0 !important;
            }

            table[class="body"] .btn table {
                width: 100% !important;
            }

            table[class="body"] .btn a {
                width: 100% !important;
            }

            table[class="body"] .img-responsive {
                height: auto !important;
                max-width: 100% !important;
                width: auto !important;
            }
        }

        body {
            background: linear-gradient(135deg, #5c9dff, #d08bff);
            font-family: 'Arial', sans-serif;
            font-size: 14px;
            line-height: 1.4;
            margin: 0;
            padding: 0;
            -webkit-font-smoothing: antialiased;
            -webkit-text-size-adjust: 100%;
        }

        .container {
            max-width: 580px;
            margin: 20px auto;
            padding: 20px;
        }

        .main {
            background: #ffffff;
            border-radius: 10px;
            overflow: hidden;
            box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
            padding: 20px;
        }

        .header img {
            max-width: 100px;
        }

        .title {
            font-size: 24px;
            color: #2c3e50;
            text-align: center;
            margin-bottom: 20px;
        }

        .content {
            color: #34495e;
            font-size: 16px;
            text-align: center;
        }

        .code {
            font-size: 36px;
            font-weight: bold;
            color: #5c9dff;
            background: #eaf3ff;
            padding: 10px;
            border-radius: 5px;
            display: inline-block;
            margin: 20px 0;
        }

        .footer {
            color: #7f8c8d;
            font-size: 12px;
            text-align: center;
            margin-top: 20px;
        }
    </style>
</head>
<body>
<table role="presentation" border="0" cellpadding="0" cellspacing="0" class="body" width="100%">
    <tr>
        <td align="center">
            <div class="container">
                <div class="main">
                    <h1 class="title">Подтверждение регистрации</h1>
                    <p class="content">Добро пожаловать в LEARNIO!<br>Ваш код активации:</p>
                    <p class="code">880321</p>
                    <p class="content">
                        Введите этот код для завершения регистрации. Если вы не запрашивали код, просто проигнорируйте это сообщение.
                    </p>
                    <div class="footer">© 2024 LEARNIO. Все права защищены.</div>
                </div>
            </div>
        </td>
    </tr>
</table>
</body>
</html>`

		content, err = changeHTMLValue(content, "code", code)
		if err != nil {
			logger.Log.Error(err)
			return
		}

		if err := u.mailer.SendEmail(models.Gmail{
			Subject: "Подтверждение пароля",
			Content: content,
			TO:      []string{request.Email},
		}); err != nil {
			logger.Log.Errorf("send email, err: %v", err)
		}

		logger.Log.Info("send email")
	}(code)

	return id, nil
}

func (u *UCUser) ConfirmMail(ctx context.Context, userID int, code string) error {
	userCode := u.redis.GetUsersRegistrationCode(ctx, userID)

	if strings.EqualFold(code, userCode) {
		return errors.New("invalid code")
	}

	if err := u.user.ConfirmEmail(ctx, userID); err != nil {
		return errors.Wrap(err, "confirm email")
	}

	return nil
}

func (u *UCUser) Login(ctx context.Context, login, password string) (*models.User, error) {
	user, err := u.user.SelectUserByLogin(ctx, login)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	if !user.ConfirmEmail {
		return nil, errors.New("user not confirm email")
	}

	err = u.compareHashAndPassword(user.HashPass, password)
	if err != nil {
		return nil, errors.New("wrong password or email")
	}

	return user, nil
}

func (u *UCUser) createHashPassword(password string) (string, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("create hashed password was failed: %v", err.Error())
	}

	return string(hashPassword), nil
}

func (u *UCUser) compareHashAndPassword(hash, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return fmt.Errorf("compare is wrong, err: %w", err)
	}

	return nil
}

func changeHTMLValue(email string, class string, newValue string) (string, error) {
	// Используем goquery для парсинга HTML-контента из поля Body
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(email))
	if err != nil {
		return "", errors.Wrap(err, "can not parse html content")
	}

	// Находим элементы с указанным классом
	doc.Find("." + class).Each(func(i int, s *goquery.Selection) {
		// Меняем значение внутри найденных элементов
		s.SetHtml(newValue)
	})

	// Обновляем значение в поле Body структуры Email
	htmlContent, err := doc.Html()
	if err != nil {
		return "", errors.Wrap(err, "can not get html content")
	}

	return htmlContent, nil
}

func generateVerificationCode() string {
	b := make([]byte, entity.LenRegistrationCode)

	_, _ = rand.Read(b)

	code := make([]rune, entity.LenRegistrationCode)
	for i := range entity.LenRegistrationCode {
		index := int(b[i]) % len(entity.AllSymbol)
		code[i] = entity.AllSymbol[index]
	}

	return string(code)
}
