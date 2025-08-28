package usecases

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/dto"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/styles"
	"github.com/google/uuid"
	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
	"go.uber.org/zap"
)

// ModernPDFUseCase defines the interface for modern PDF operations
type ModernPDFUseCase interface {
	GeneratePDF(ctx context.Context, curriculumID uuid.UUID) error
}

// ModernPDFUseCaseImpl is responsible for generating PDFs of curriculums
type ModernPDFUseCaseImpl struct {
	curriculumUseCase CurriculumUseCase
	logger            *zap.Logger
}

// NewModernPDFUseCase creates a new instance of the PDF generator
func NewModernPDFUseCase(curriculumUseCase CurriculumUseCase, logger *zap.Logger) ModernPDFUseCase {
	return &ModernPDFUseCaseImpl{
		curriculumUseCase: curriculumUseCase,
		logger:            logger,
	}
}

// GeneratePDF generates the PDF of the curriculum and saves it in the specified paths
func (pg *ModernPDFUseCaseImpl) GeneratePDF(ctx context.Context, curriculumID uuid.UUID) error {
	pg.logger.Info("Starting PDF generation",
		zap.String("curriculum_id", curriculumID.String()),
	)

	// Get the specific curriculum from the database
	curriculum, err := pg.curriculumUseCase.GetCurriculumByID(ctx, curriculumID)
	if err != nil {
		pg.logger.Error("Failed to get curriculum from database",
			zap.String("curriculum_id", curriculumID.String()),
			zap.Error(err),
		)
		return err
	}

	pg.logger.Info("Curriculum retrieved successfully",
		zap.String("curriculum_id", curriculumID.String()),
		zap.String("full_name", curriculum.FullName),
	)

	maroto := pg.createMaroto(*curriculum)

	pg.logger.Debug("Maroto document created, generating PDF",
		zap.String("curriculum_id", curriculumID.String()),
	)

	document, err := maroto.Generate()
	if err != nil {
		pg.logger.Error("Failed to generate PDF document",
			zap.String("curriculum_id", curriculumID.String()),
			zap.Error(err),
		)
		return err
	}

	pg.logger.Debug("PDF document generated successfully",
		zap.String("curriculum_id", curriculumID.String()),
	)

	// Create directories if they don't exist
	pdfDir := os.Getenv("PDF_DIR")
	textDir := os.Getenv("TEXT_DIR")

	if err := pg.ensureDirectoryExists(pdfDir); err != nil {
		pg.logger.Error("Failed to create PDF directory",
			zap.String("directory", pdfDir),
			zap.String("curriculum_id", curriculumID.String()),
			zap.Error(err),
		)
		return err
	}
	if err := pg.ensureDirectoryExists(textDir); err != nil {
		pg.logger.Error("Failed to create text directory",
			zap.String("directory", textDir),
			zap.String("curriculum_id", curriculumID.String()),
			zap.Error(err),
		)
		return err
	}

	// Save the PDF
	pdfPath := filepath.Join(pdfDir, "curriculum_modern.pdf")
	err = document.Save(pdfPath)
	if err != nil {
		pg.logger.Error("Failed to save PDF file",
			zap.String("file_path", pdfPath),
			zap.String("curriculum_id", curriculumID.String()),
			zap.Error(err),
		)
		return err
	}

	pg.logger.Info("PDF file saved successfully",
		zap.String("file_path", pdfPath),
		zap.String("curriculum_id", curriculumID.String()),
	)

	// Save the report in text
	textPath := filepath.Join(textDir, "curriculum_modern_report.txt")
	err = document.GetReport().Save(textPath)
	if err != nil {
		pg.logger.Error("Failed to save text report",
			zap.String("file_path", textPath),
			zap.String("curriculum_id", curriculumID.String()),
			zap.Error(err),
		)
		return err
	}

	pg.logger.Info("Text report saved successfully",
		zap.String("file_path", textPath),
		zap.String("curriculum_id", curriculumID.String()),
	)

	pg.logger.Info("PDF generation completed successfully",
		zap.String("curriculum_id", curriculumID.String()),
		zap.String("pdf_path", pdfPath),
		zap.String("text_path", textPath),
	)

	return nil
}

// ensureDirectoryExists creates the directory if it doesn't exist
func (pg *ModernPDFUseCaseImpl) ensureDirectoryExists(dir string) error {
	pg.logger.Debug("Ensuring directory exists",
		zap.String("directory", dir),
	)
	return os.MkdirAll(dir, 0755)
}

// createMaroto creates and configures the Maroto document
func (pg *ModernPDFUseCaseImpl) createMaroto(curriculum dto.CurriculumResponse) core.Maroto {
	cfg := config.NewBuilder().
		WithPageNumber().
		WithLeftMargin(25).
		WithTopMargin(25).
		WithRightMargin(25).
		WithBottomMargin(20).
		Build()

	mrt := maroto.New(cfg)
	m := maroto.NewMetricsDecorator(mrt)

	// Add all sections
	m.AddRows(pg.getHeaderSection(curriculum)...)
	m.AddRows(pg.getAboutSection(curriculum)...)
	m.AddRows(pg.getExperienceSection(curriculum.Works)...)
	m.AddRows(pg.getSkillsSection(curriculum)...)
	m.AddRows(pg.getEducationSection(curriculum)...)
	m.AddRows(pg.getAdditionalInfoSection(curriculum)...)

	return m
}

// getHeaderSection returns the elegant header section
func (pg *ModernPDFUseCaseImpl) getHeaderSection(curriculum dto.CurriculumResponse) []core.Row {
	return []core.Row{
		// Name
		row.New(20).Add(
			text.NewCol(12, curriculum.FullName, props.Text{
				Top:   5,
				Size:  36,
				Style: fontstyle.Bold,
				Align: align.Left,
				Color: styles.GetPrimaryColorModernPDF(),
			}),
		),

		// Job title
		row.New(12).Add(
			text.NewCol(12, curriculum.JobDescription, props.Text{
				Top:   0,
				Size:  20,
				Style: fontstyle.Normal,
				Align: align.Left,
				Color: styles.GetSecondaryColorModernPDF(),
			}),
		),

		// Contact info
		row.New(10).Add(
			text.NewCol(6, curriculum.Email, props.Text{
				Top:   10,
				Size:  11,
				Align: align.Left,
				Color: styles.GetMutedTextColorModernPDF(),
			}),
			text.NewCol(6, curriculum.Phone, props.Text{
				Top:   10,
				Size:  11,
				Align: align.Right,
				Color: styles.GetMutedTextColorModernPDF(),
			}),
		),

		// Divider line
		row.New(15).Add(
			text.NewCol(12, "________________________________________________________________________________", props.Text{
				Top:   8,
				Size:  10,
				Align: align.Center,
				Color: styles.GetBorderColorModernPDF(),
			}),
		),

		row.New(10),
	}
}

// getAboutSection returns the about section
func (pg *ModernPDFUseCaseImpl) getAboutSection(curriculum dto.CurriculumResponse) []core.Row {
	return []core.Row{
		// Section title
		row.New(12).Add(
			text.NewCol(12, "ABOUT", props.Text{
				Top:   3,
				Size:  16,
				Style: fontstyle.Bold,
				Align: align.Left,
				Color: styles.GetPrimaryColorModernPDF(),
			}),
		),

		// About content
		row.New(15).Add(
			text.NewCol(12, curriculum.Intro, props.Text{
				Top:   2,
				Size:  12,
				Align: align.Justify,
				Color: styles.GetTextColorModernPDF(),
			}),
		),

		row.New(12),
	}
}

// getExperienceSection returns the experience section
func (pg *ModernPDFUseCaseImpl) getExperienceSection(works []dto.WorkResponse) []core.Row {
	var rows []core.Row

	// Section title
	rows = append(rows, row.New(12).Add(
		text.NewCol(12, "PROFESSIONAL EXPERIENCE", props.Text{
			Top:   3,
			Size:  16,
			Style: fontstyle.Bold,
			Align: align.Left,
			Color: styles.GetPrimaryColorModernPDF(),
		}),
	))

	// Experience items
	for i, work := range works {
		rows = append(rows, pg.getExperienceItem(work)...)

		// Add spacing between items (except for the last one)
		if i < len(works)-1 {
			rows = append(rows, row.New(8))
		}
	}

	rows = append(rows, row.New(12))
	return rows
}

// getExperienceItem returns a single experience item
func (pg *ModernPDFUseCaseImpl) getExperienceItem(work dto.WorkResponse) []core.Row {
	var rows []core.Row

	endDateStr := "Present"
	if work.EndDate != nil {
		endDateStr = work.EndDate.Format("Jan 2006")
	}

	// Job title and period
	rows = append(rows, row.New(8).Add(
		text.NewCol(8, work.JobTitle, props.Text{
			Top:   2,
			Size:  14,
			Style: fontstyle.Bold,
			Align: align.Left,
			Color: styles.GetSecondaryColorModernPDF(),
		}),
		text.NewCol(4, work.StartDate.Format("Jan 2006")+" - "+endDateStr, props.Text{
			Top:   2,
			Size:  11,
			Style: fontstyle.Italic,
			Align: align.Right,
			Color: styles.GetMutedTextColorModernPDF(),
		}),
	))

	// Company name
	rows = append(rows, row.New(6).Add(
		text.NewCol(12, work.CompanyName, props.Text{
			Top:   1,
			Size:  13,
			Style: fontstyle.Bold,
			Align: align.Left,
			Color: styles.GetPrimaryColorModernPDF(),
		}),
	))

	// Company description
	rows = append(rows, row.New(12).Add(
		text.NewCol(12, work.CompanyDescription, props.Text{
			Top:   2,
			Size:  11,
			Align: align.Justify,
			Color: styles.GetTextColorModernPDF(),
		}),
	))

	return rows
}

// getSkillsSection returns the skills section
func (pg *ModernPDFUseCaseImpl) getSkillsSection(curriculum dto.CurriculumResponse) []core.Row {
	var rows []core.Row

	// Section title
	rows = append(rows, row.New(12).Add(
		text.NewCol(12, "SKILLS & TECHNOLOGIES", props.Text{
			Top:   3,
			Size:  16,
			Style: fontstyle.Bold,
			Align: align.Left,
			Color: styles.GetPrimaryColorModernPDF(),
		}),
	))

	// Technologies
	technologies := strings.Split(curriculum.Technologies, ",")
	var techList []string
	for _, tech := range technologies {
		tech = strings.TrimSpace(tech)
		if tech != "" {
			techList = append(techList, "• "+tech)
		}
	}

	// Split technologies into columns for better layout
	techText := strings.Join(techList, "    ")
	rows = append(rows, row.New(12).Add(
		text.NewCol(12, techText, props.Text{
			Top:   2,
			Size:  11,
			Align: align.Left,
			Color: styles.GetTextColorModernPDF(),
		}),
	))

	rows = append(rows, row.New(12))
	return rows
}

// getEducationSection returns the education section
func (pg *ModernPDFUseCaseImpl) getEducationSection(curriculum dto.CurriculumResponse) []core.Row {
	return []core.Row{
		// Section title
		row.New(12).Add(
			text.NewCol(12, "EDUCATION", props.Text{
				Top:   3,
				Size:  16,
				Style: fontstyle.Bold,
				Align: align.Left,
				Color: styles.GetPrimaryColorModernPDF(),
			}),
		),

		// Education content
		row.New(10).Add(
			text.NewCol(12, curriculum.LevelEducation, props.Text{
				Top:   2,
				Size:  12,
				Align: align.Left,
				Color: styles.GetTextColorModernPDF(),
			}),
		),

		row.New(12),
	}
}

// getAdditionalInfoSection returns the additional information section
func (pg *ModernPDFUseCaseImpl) getAdditionalInfoSection(curriculum dto.CurriculumResponse) []core.Row {
	var rows []core.Row

	// Section title
	rows = append(rows, row.New(12).Add(
		text.NewCol(12, "ADDITIONAL INFORMATION", props.Text{
			Top:   3,
			Size:  16,
			Style: fontstyle.Bold,
			Align: align.Left,
			Color: styles.GetPrimaryColorModernPDF(),
		}),
	))

	// Languages
	rows = append(rows, row.New(8).Add(
		text.NewCol(4, "Languages:", props.Text{
			Top:   2,
			Size:  12,
			Style: fontstyle.Bold,
			Align: align.Left,
			Color: styles.GetSecondaryColorModernPDF(),
		}),
		text.NewCol(8, curriculum.Languages, props.Text{
			Top:   2,
			Size:  11,
			Align: align.Left,
			Color: styles.GetTextColorModernPDF(),
		}),
	))

	// Driver's License
	rows = append(rows, row.New(8).Add(
		text.NewCol(4, "Driver's License:", props.Text{
			Top:   2,
			Size:  12,
			Style: fontstyle.Bold,
			Align: align.Left,
			Color: styles.GetSecondaryColorModernPDF(),
		}),
		text.NewCol(8, curriculum.DriverLicense, props.Text{
			Top:   2,
			Size:  11,
			Align: align.Left,
			Color: styles.GetTextColorModernPDF(),
		}),
	))

	// Courses
	rows = append(rows, row.New(12).Add(
		text.NewCol(4, "Courses:", props.Text{
			Top:   2,
			Size:  12,
			Style: fontstyle.Bold,
			Align: align.Left,
			Color: styles.GetSecondaryColorModernPDF(),
		}),
		text.NewCol(8, curriculum.Courses, props.Text{
			Top:   2,
			Size:  11,
			Align: align.Left,
			Color: styles.GetTextColorModernPDF(),
		}),
	))

	// Social Links
	rows = append(rows, row.New(12).Add(
		text.NewCol(4, "Social Links:", props.Text{
			Top:   2,
			Size:  12,
			Style: fontstyle.Bold,
			Align: align.Left,
			Color: styles.GetSecondaryColorModernPDF(),
		}),
		text.NewCol(8, curriculum.SocialLinks, props.Text{
			Top:   2,
			Size:  11,
			Align: align.Left,
			Color: styles.GetTextColorModernPDF(),
		}),
	))

	rows = append(rows, row.New(15))
	return rows
}
