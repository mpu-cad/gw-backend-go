package user

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/mpu-cad/gw-backend-go/internal/logger"

	"golang.org/x/crypto/bcrypt"

	"github.com/mpu-cad/gw-backend-go/internal/models"
)

type UCUser struct {
	user   userRepos
	mailer mailer
}

func NewUCUser(user userRepos, mailer mailer) *UCUser {
	return &UCUser{
		user:   user,
		mailer: mailer,
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

	go func() {
		code := generateVerificationCode()

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

		err, content = changeHTMLValue(content, "code", code)

		if err := u.mailer.SendEmail(models.Gmail{
			Subject: "Подтверждение пароля",
			Content: content,
			TO:      []string{request.Email},
		}); err != nil {
			logger.Log.Errorf("send email, err: %v", err)
		}

		logger.Log.Info("send email")
	}()

	return id, nil
}

func (u *UCUser) Login(ctx context.Context, email, password string) (*models.User, error) {
	user, err := u.user.SelectUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
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

func changeHTMLValue(email string, class string, newValue string) (error, string) {
	// Используем goquery для парсинга HTML-контента из поля Body
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(email))
	if err != nil {
		return err, ""
	}

	// Находим элементы с указанным классом
	doc.Find("." + class).Each(func(i int, s *goquery.Selection) {
		// Меняем значение внутри найденных элементов
		s.SetHtml(newValue)
	})

	// Обновляем значение в поле Body структуры Email
	htmlContent, err := doc.Html()
	if err != nil {
		return err, ""
	}

	return nil, htmlContent
}

func generateVerificationCode() string {
	code := make([]string, 6)
	for i := 0; i < 6; i++ {
		if rand.Intn(2) == 1 {
			code[i] = strconv.Itoa(rand.Intn(10))
		} else {
			code[i] = string(rune(rand.Intn(128)%26 + 65))
		}
	}
	return strings.Join(code, "")
}
